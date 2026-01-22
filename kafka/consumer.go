package kafka

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"

	"github.com/agurinov/gopl/backoff"
	"github.com/agurinov/gopl/graceful"
	c "github.com/agurinov/gopl/patterns/creational"
	"github.com/agurinov/gopl/run"
)

type (
	Consumer interface {
		Run(context.Context) error
		Close(context.Context) error
	}
	consumer[R any, V any] struct {
		metrics         consumerMetrics
		recordDiscarder RecordDiscarder[V]
		recordMapper    RecordMapper[R, V]
		logger          *zap.Logger
		partitionHolder *partitionHolder
		clientMu        *sync.Mutex
		client          *kgo.Client
		handler         Handler[R]
		handlerBatch    HandlerBatch[R]
		config          config
		graceful.Wrapper
	}
)

type (
	ConsumerOption[R any] c.Option[kgoConsumer[R]]
)

func (cs consumer[R, V]) processPartition(
	ctx context.Context,
	partition int32,
) error {
	l := cs.logger.With(
		zap.String("eventloop.topic", cs.config.topic),
		zap.Int32("eventloop.partition", partition),
	)

	l.Info("started partition eventloop")
	defer l.Info("exited partition eventloop")

	bopts := cs.config.idleBackoffOptions(l)

	idleBackoff, err := backoff.New(bopts...)
	if err != nil {
		return fmt.Errorf("can't create eventloop backoff: %w", err)
	}

	iterationFn := run.ErrorFn(func() error {
		l.Info(
			"polling partition",
			zap.Int64("buffered.count", cs.client.BufferedFetchRecords()),
		)

		fetches := cs.pollPartition(ctx, partition)

		switch cs.analyzeFetches(fetches, l) {
		case processAction:
			// move straight
		case exitAction:
			return nil
		case skipAction:
			continue
		case waitAction:
			if _, bErr := idleBackoff.Wait(ctx); bErr != nil {
				return fmt.Errorf("can't wait idle backoff: %w", bErr)
			}

			continue
		}
	})

	// TODO: graceful.Run ?
	for {
		l.Info(
			"polling partition",
			zap.Int64("buffered.count", cs.client.BufferedFetchRecords()),
		)

		fetches := cs.pollPartition(ctx, partition)

		switch cs.analyzeFetches(fetches, l) {
		case processAction:
			// move straight
		case exitAction:
			return nil
		case skipAction:
			continue
		case waitAction:
			if _, bErr := idleBackoff.Wait(ctx); bErr != nil {
				return fmt.Errorf("can't wait idle backoff: %w", bErr)
			}

			continue
		}

		switch {
		case cs.handlerBatch != nil:
			fetches.EachPartition(
				cs.eachBatchFunc(ctx, partition),
			)
		case cs.handler != nil:
			fetches.EachRecord(
				cs.eachRecordFunc(ctx, partition),
			)
		}

		// catch via channel?

		if dErr := cs.recordDiscarder.Discard(ctx); dErr != nil {
			return fmt.Errorf("can't discard records: %w", dErr)
		}

		// TODO: commit / dlq - discarder
		// NOTE: we should develop a decision about batch results. Proposals:
		// 1) Client handler return `all or nothing`
		// 2) Client handler with get a struct which is buffers and can dedicate records to dlq or
		// commit
		_ = cs.metrics
	}
}

func (cs consumer[R, V]) Close() {
	fn := run.SimpleFn(func() {
		cs.client.Close()
	})

	cs.Wrapper.Close(fn)(nil)
}

func (cs consumer[R, V]) Run() error {
	fn := run.SimpleFn(func() {
		time.Sleep(20 * time.Second)

		cs.logger.Info(
			"consumer ping",
			zap.String("topic", cs.config.topic),
			zap.Int32s("assigned.partitions", cs.partitionHolder.assignedPartitions()),
		)
	})

	return cs.Wrapper.Run(fn)(nil)
}

// NOTE: main logic appears here.
// According to kafka's documentation we CAN use pause/resume logic while we optimize goroutine
// consurrency.
// Proofs:
//  1. https://docs.confluent.io/platform/current/clients/javadocs/javadoc/org/apache/kafka/clients/consumer/KafkaConsumer.html#consumption-flow-control-heading
//  2. https://docs.confluent.io/platform/current/clients/javadocs/javadoc/org/apache/kafka/clients/consumer/KafkaConsumer.html#1-one-consumer-per-thread-heading
//
// franz-go's logic of polling separated into 2 steps:
//  1. fetch. Single eventloop across all assigned partitions with buffering feature. It is hidden and immutable for us.
//  2. poll. We can switch client to desired partition to fetch only those records. While we can't
//     fully clone the client mutex appears here only on roundtrip.
//
// There are cases when partition doesn't have records ahead and this Poll invoke can hung.
// We've limited roundtrip from config Duration and this case (deadline exceeded) for NOW interprets as IDLE state.
func (cs consumer[R, V]) pollPartition(
	ctx context.Context,
	partition int32,
) kgo.Fetches {
	forPause := map[string][]int32{
		cs.config.topic: cs.partitionHolder.assignedPartitions(),
	}

	forResume := map[string][]int32{
		cs.config.topic: {partition},
	}

	cs.clientMu.Lock()
	defer cs.clientMu.Unlock()

	cs.client.PauseFetchPartitions(forPause)
	cs.client.ResumeFetchPartitions(forResume)

	defer cs.client.AllowRebalance()

	ctx, cancel := context.WithTimeout(ctx, cs.config.maxPollDuration)
	defer cancel()

	return cs.client.PollRecords(ctx, cs.config.maxPollRecords)
}

func NewConsumer[R any](opts ...ConsumerOption[R]) (Consumer, error) {
	obj := kgoConsumer[R]{
		clientMu: new(sync.Mutex),
		config: config{
			idleBackoffMinDelay: 100 * time.Millisecond,
			idleBackoffMaxDelay: 3 * time.Second,
			maxPollRecords:      12,
			maxPollDuration:     200 * time.Millisecond,
		},
		partitionHolder: &partitionHolder{
			assigned: make(map[int32]partitionContext, 12),
		},
		recordMapper:    kgoRecordMapper[R]{},
		recordDiscarder: noopRecordDiscarder{},
	}

	obj, err := c.ConstructWithValidate(obj, opts...)
	if err != nil {
		return nil, err
	}

	gracefulWrapper, err := graceful.NewWrapper(
		graceful.WithWrapperLogger(obj.logger),
	)
	if err != nil {
		return nil, err
	}

	obj.Wrapper = gracefulWrapper

	return obj, nil
}
