package kafka

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"

	"github.com/agurinov/gopl/backoff"
	c "github.com/agurinov/gopl/patterns/creational"
)

type (
	Consumer interface {
		Consume(context.Context) error
		Close(context.Context) error
	}
	consumer struct {
		metrics         consumerMetrics
		recordDiscarder RecordDiscarder[*kgo.Record]
		recordMapper    RecordMapper[*kgo.Record, Record]
		logger          *zap.Logger
		closed          *atomic.Bool
		partitionHolder *partitionHolder
		clientMu        *sync.Mutex
		client          *kgo.Client
		handler         Handler
		handlerBatch    HandlerBatch
		config          config
	}
	ConsumerOption c.Option[consumer]
)

func (c *consumer) processPartition(
	ctx context.Context,
	partition int32,
) {
	l := c.logger.With(
		zap.String("eventloop.topic", c.config.topic),
		zap.Int32("eventloop.partition", partition),
	)

	l.Info("started partition eventloop")
	defer l.Info("exited partition eventloop")

	bopts := c.config.idleBackoffOptions(l)

	idleBackoff, err := backoff.New(bopts...)
	if err != nil {
		l.Error(
			"can't create eventloop backoff",
			zap.Error(err),
		)

		return
	}

	for {
		l.Info(
			"polling partition",
			zap.Int64("buffered.count", c.client.BufferedFetchRecords()),
		)

		fetches := c.pollPartition(ctx, partition)

		switch c.analyzeFetches(fetches, l) {
		case processAction:
			// move straight
		case exitAction:
			return
		case skipAction:
			continue
		case waitAction:
			if _, bErr := idleBackoff.Wait(ctx); bErr != nil {
				l.Error(
					"can't wait idle backoff",
					zap.Error(bErr),
				)

				return
			}

			continue
		}

		switch {
		case c.handlerBatch != nil:
			fetches.EachPartition(
				c.eachBatchFunc(ctx, partition),
			)
		case c.handler != nil:
			fetches.EachRecord(
				c.eachRecordFunc(ctx, partition),
			)
		}

		// catch via channel?

		if err := c.recordDiscarder.Discard(ctx); err != nil {
			l.Error(
				"can't discard records",
				zap.Error(err),
			)

			return
		}

		// TODO: commit / dlq - discarder
		// NOTE: we should develop a decision about batch results. Proposals:
		// 1) Client handler return `all or nothing`
		// 2) Client handler with get a struct which is buffers and can dedicate records to dlq or
		// commit
		_ = c.metrics
	}
}

// Embed graceful ??
func (c *consumer) Close() {
	if !c.closed.CompareAndSwap(false, true) {
		return
	}

	c.logger.Info("closing Kafka Consumer")

	for _, p := range c.partitionHolder.assignedPartitions() {
		c.partitionHolder.revokePartition(p)
	}

	c.client.Close()
}

func (c *consumer) IsClosed(ctxs ...context.Context) bool {
	if c.closed.Load() {
		return true
	}

	ctxs = append(ctxs, c.client.Context())

	for i := range ctxs {
		if ctxs[i] == nil {
			continue
		}

		select {
		case <-ctxs[i].Done():
			return true
		default:
		}
	}

	return false
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
func (c *consumer) pollPartition(
	ctx context.Context,
	partition int32,
) kgo.Fetches {
	forPause := map[string][]int32{
		c.config.topic: c.partitionHolder.assignedPartitions(),
	}

	forResume := map[string][]int32{
		c.config.topic: {partition},
	}

	c.clientMu.Lock()
	defer c.clientMu.Unlock()

	c.client.PauseFetchPartitions(forPause)
	c.client.ResumeFetchPartitions(forResume)

	defer c.client.AllowRebalance()

	ctx, cancel := context.WithTimeout(ctx, c.config.maxPollDuration)
	defer cancel()

	return c.client.PollRecords(ctx, c.config.maxPollRecords)
}

func NewConsumer(opts ...ConsumerOption) (consumer, error) {
	obj := consumer{
		closed:   new(atomic.Bool),
		clientMu: new(sync.Mutex),
		config: config{
			idleBackoffMinDelay: 100 * time.Millisecond,
			idleBackoffMaxDelay: 3 * time.Second,
			maxPollRecords:      12,
			maxPollDuration:     200 * time.Millisecond,
		},
		partitionHolder: &partitionHolder{
			assigned: map[int32]partitionContext{},
		},
		recordMapper:    kgoRecordMapper{},
		recordDiscarder: noopRecordDiscarder{},
	}

	return c.ConstructWithValidate(obj, opts...)
}
