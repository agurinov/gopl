package kafka

import (
	"context"
	"fmt"

	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"

	"github.com/agurinov/gopl/backoff"
	"github.com/agurinov/gopl/graceful"
	"github.com/agurinov/gopl/run"
)

func (cs consumer) partitionLoop(partition int32) run.Fn {
	return func(ctx context.Context) error {
		l := cs.logger.With(
			zap.String("eventloop.topic", cs.topic),
			zap.Int32("eventloop.partition", partition),
		)

		l.Info("started partition eventloop")
		defer l.Info("exited partition eventloop")

		idleBackoff, err := cs.backoffFabric.New(
			backoff.WithLogger(l),
		)
		if err != nil {
			return fmt.Errorf("can't create eventloop backoff: %w", err)
		}

		var (
			iterationFn = cs.partitionIteration(partition, idleBackoff)
			loopFn      = graceful.RunLoop(iterationFn)
		)

		return loopFn(ctx)
	}
}

func (cs consumer) partitionIteration(
	partition int32,
	idleBackoff backoff.Backoff, // TODO: check its NOT ptr
) run.Fn {
	return func(ctx context.Context) error {
		l := cs.logger.With(
			zap.String("eventloop.topic", cs.topic),
			zap.Int32("eventloop.partition", partition),
		)

		l.Info("partition polling iteration")

		fetches := cs.pollPartition(ctx, partition)

		switch cs.analyzeFetches(fetches, l) {
		case exitAction:
			return graceful.ErrStopLoop
		case skipAction:
			return nil
		case waitAction:
			if _, bErr := idleBackoff.Wait(ctx); bErr != nil {
				return fmt.Errorf("can't wait idle backoff: %w", bErr)
			}

			return nil
		case processAction:
			idleBackoff.Reset()
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

		// TODO: catch via channel?
		// NOTE: we should develop a decision about batch results. Proposals:
		// 1) Client handler return `all or nothing`
		// 2) Client handler with get a struct which is buffers and can dedicate records to dlq or
		var (
			recordsToDiscard []Record
			recordsToCommit  []Record
		)

		if dErr := cs.recordDiscarder.Discard(ctx, recordsToDiscard...); dErr != nil {
			return fmt.Errorf("can't discard records: %w", dErr)
		}

		if cErr := cs.client.CommitRecords(ctx, recordsToCommit...); cErr != nil {
			return fmt.Errorf("can't commit records: %w", cErr)
		}

		return nil
	}
}

func (cs consumer) pollPartition(
	ctx context.Context,
	partition int32,
) kgo.Fetches {
	ctx, cancel := context.WithTimeout(ctx, cs.maxPollDuration)
	defer cancel()

	allPartitions := cs.partitionDispatcher.Running()

	return cs.client.PollTopicPartition(
		ctx,
		cs.topic,
		partition,
		allPartitions,
		cs.maxPollRecords,
	)
}
