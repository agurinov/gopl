package kafka

import (
	"cmp"
	"context"

	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"

	"github.com/agurinov/gopl/x"
)

type (
	Handler      func(context.Context, Record) error
	HandlerBatch func(context.Context, []Record) error
)

func (c *consumer) eachRecordFunc(
	ctx context.Context,
	partition int32,
) func(*kgo.Record) {
	l := c.logger.With(
		zap.String("eventloop.topic", c.config.topic),
		zap.Int32("eventloop.partition", partition),
	)

	return func(r *kgo.Record) {
		isNotMyTopicPartition := cmp.Or(
			r.Topic != c.config.topic,
			r.Partition != partition,
		)

		if isNotMyTopicPartition {
			l.Warn(
				"not my topic partition; skipping",
				zap.String("record.topic", r.Topic),
				zap.Int32("record.partition", r.Partition),
			)

			return
		}

		record := recordFromKgo(r)

		if err := c.handler(ctx, record); err != nil {
			l.Error(
				"can't handle record",
				zap.Error(err),
			)
		}
	}
}

func (c *consumer) eachBatchFunc(
	ctx context.Context,
	partition int32,
) func(tp kgo.FetchTopicPartition) {
	l := c.logger.With(
		zap.String("eventloop.topic", c.config.topic),
		zap.Int32("eventloop.partition", partition),
	)

	return func(tp kgo.FetchTopicPartition) {
		isNotMyTopicPartition := cmp.Or(
			tp.Topic != c.config.topic,
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
			recordFromKgo,
		)

		if err := c.handlerBatch(ctx, records); err != nil {
			l.Error(
				"can't handle batch of records",
				zap.Int("batch.size", len(tp.Records)),
				zap.Int64("offset.latest", x.Last(tp.Records).Offset),
				zap.Error(err),
			)
		}
	}
}
