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

			return
		}

		record := cs.recordMapper.FromVendor(r)

		if err := cs.handler(ctx, record); err != nil {
			l.Error(
				"can't handle record",
				zap.Error(err),
			)
		}
	}
}

func (cs consumer) eachBatchFunc(
	ctx context.Context,
	partition int32,
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

			return
		}

		records := x.SliceConvert(
			tp.Records,
			cs.recordMapper.FromVendor,
		)

		if err := cs.handlerBatch(ctx, records); err != nil {
			l.Error(
				"can't handle batch of records",
				zap.Int("batch.size", len(records)),
				zap.Int64("offset.latest", x.Last(records).Offset),
				zap.Error(err),
			)
		}
	}
}
