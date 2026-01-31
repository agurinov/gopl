package kafka

import (
	"cmp"
	"context"

	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"

	"github.com/agurinov/gopl/x"
)

type (
	Handler      = func(context.Context, Record) error
	HandlerBatch = func(context.Context, []Record) error
)

func (cs consumer) eachRecordFunc(
	ctx context.Context,
	partition int32,
	discardCh chan<- Record,
	commitCh chan<- Record,
) func(Record) {
	l := cs.logger.With(
		zap.String("eventloop.topic", cs.topic),
		zap.Int32("eventloop.partition", partition),
	)

	return func(r Record) {
		isNotMyTopicPartition := cmp.Or(
			r.Topic != cs.topic,
			r.Partition != partition,
		)
		if isNotMyTopicPartition {
			l.Warn(
				"not my topic partition; skipping",
				RecordLogFields(r)...,
			)

			cs.metrics.notMyRecordInc(r.Topic, r.Partition)

			return
		}

		record := cs.recordMapper.FromVendor(r)

		if err := cs.handler(ctx, record); err != nil {
			l.Error(
				"can't handle record",
				zap.Error(err),
			)

			discardCh <- record
		}

		cs.hooks.partitionCounters.Inc(partition)

		commitCh <- record
	}
}

func (cs consumer) eachBatchFunc(
	ctx context.Context,
	partition int32,
	discardCh chan<- Record,
	commitCh chan<- Record,
) func(tp kgo.FetchTopicPartition) {
	l := cs.logger.With(
		zap.String("eventloop.topic", cs.topic),
		zap.Int32("eventloop.partition", partition),
	)

	return func(tp kgo.FetchTopicPartition) {
		isNotMyTopicPartition := cmp.Or(
			tp.Topic != cs.topic,
			tp.Partition != partition,
		)
		if isNotMyTopicPartition {
			l.Warn(
				"not my topic partition; skipping",
				zap.String("batch.topic", tp.Topic),
				zap.Int32("batch.partition", tp.Partition),
			)

			cs.metrics.notMyRecordInc(tp.Topic, tp.Partition)

			return
		}

		records := x.SliceConvert(
			tp.Records,
			cs.recordMapper.FromVendor,
		)

		// Strategy for now: all or nothing.
		batchErr := cs.handlerBatch(ctx, records)

		if batchErr != nil {
			l.Error(
				"can't handle batch of records",
				zap.Int("batch.size", len(records)),
				zap.Int64("offset.latest", x.Last(records).Offset),
				zap.Error(batchErr),
			)
		}

		for _, r := range records {
			if batchErr != nil {
				discardCh <- r
			}

			cs.hooks.partitionCounters.Inc(partition)

			commitCh <- r
		}
	}
}
