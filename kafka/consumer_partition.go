package kafka

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"

	"github.com/agurinov/gopl/backoff"
	"github.com/agurinov/gopl/graceful"
	irun "github.com/agurinov/gopl/internal/run"
	"github.com/agurinov/gopl/run"
	"github.com/agurinov/gopl/x"
)

func (cs consumer) partitionLoop(ctx context.Context) error {
	partition, err := irun.GetDispatcherKey[int32](ctx)
	if err != nil {
		return fmt.Errorf("can't get partition key from context: %w", err)
	}

	l := cs.logger.With(
		zap.String("eventloop.topic", cs.topic),
		zap.Int32("eventloop.partition", partition),
	)

	// TODO: remove it or parametrize it. Filter out partitions
	if isSkip := slices.Contains([]int32{1, 5, 4}, partition); isSkip {
		l.Info("skip partition eventloop")

		return nil
	}

	idleBackoff, err := backoff.New(cs.backoffOptions...)
	if err != nil {
		return fmt.Errorf("can't create eventloop backoff: %w", err)
	}

	var (
		iterationFn = cs.partitionIteration(partition, idleBackoff)
		loopFn      = graceful.RunLoop(iterationFn)
	)

	l.Info("started partition eventloop")
	defer l.Info("exited partition eventloop")

	return loopFn(ctx)
}

func (cs consumer) partitionIteration(
	partition int32,
	idleBackoff backoff.Backoff,
) run.Fn {
	return func(ctx context.Context) error {
		l := cs.logger.With(
			zap.String("eventloop.topic", cs.topic),
			zap.Int32("eventloop.partition", partition),
		)

		l.Debug("partition polling iteration")

		fetches := cs.pollPartition(ctx, partition)

		switch cs.analyzeFetches(fetches, l) {
		case exitAction:
			return graceful.ErrStopLoop
		case skipAction:
			return nil
		case waitAction:
			return cs.partitionIdleWait(ctx, cs.topic, partition, idleBackoff)
		case processAction:
			idleBackoff.Reset()
		}

		var (
			discardCh = make(chan Record, cs.maxPollRecords)
			commitCh  = make(chan Record, cs.maxPollRecords)
		)

		defer func() {
			close(discardCh)
			close(commitCh)
		}()

		switch {
		case cs.handler != nil:
			fn := cs.eachRecordFunc(
				ctx,
				partition,
				discardCh,
				commitCh,
			)

			fetches.EachRecord(fn)
		case cs.handlerBatch != nil:
			fn := cs.eachBatchFunc(
				ctx,
				partition,
				discardCh,
				commitCh,
			)

			fetches.EachPartition(fn)
		}

		var (
			recordsToDiscard = x.FlattenChans(discardCh)
			recordsToCommit  = x.FlattenChans(commitCh)
		)

		if len(recordsToDiscard) != 0 {
			if dErr := cs.recordDiscarder.Discard(ctx, recordsToDiscard...); dErr != nil {
				return fmt.Errorf("can't discard records: %w", dErr)
			}

			cs.metrics.discardedAdd(
				cs.topic,
				partition,
				len(recordsToDiscard),
			)
		}

		if len(recordsToCommit) != 0 {
			if cErr := cs.client.CommitRecords(ctx, recordsToCommit...); cErr != nil {
				return fmt.Errorf("can't commit records: %w", cErr)
			}

			cs.metrics.committedAdd(
				cs.topic,
				partition,
				len(recordsToCommit),
			)
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

	allPartitions := cs.partitionDispatcher.AllKeys()

	startTime := time.Now()
	defer func() {
		cs.metrics.pollingObserve(
			cs.topic,
			partition,
			time.Since(startTime),
		)
	}()

	return cs.client.PollTopicPartition(
		ctx,
		cs.topic,
		partition,
		allPartitions,
		cs.maxPollRecords,
	)
}

func (cs consumer) partitionIdleWait(
	ctx context.Context,
	topic string,
	partition int32,
	idleBackoff backoff.Backoff,
) error {
	stat, err := idleBackoff.Wait(ctx)
	if err != nil {
		return fmt.Errorf("can't wait idle backoff: %w", err)
	}

	cs.metrics.idleObserve(
		topic,
		partition,
		stat.Duration,
	)

	return nil
}
