package kafka

import (
	"cmp"
	"context"

	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"

	"github.com/agurinov/gopl/x"
)

type (
	Handler[T any]      func(context.Context, T) error
	HandlerBatch[T any] func(context.Context, []T) error
)

func (c consumer[R, V]) eachRecordFunc(
	ctx context.Context,
	partition int32,
) func(V) {
	l := c.logger.With(
		zap.String("eventloop.topic", c.config.topic),
		zap.Int32("eventloop.partition", partition),
	)

	return func(r V) {
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

		record := c.recordMapper.FromVendor(r)

		if err := c.handler(ctx, record); err != nil {
			l.Error(
				"can't handle record",
				zap.Error(err),
			)
		}
	}
}

func (c consumer[R, V]) eachBatchFunc(
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
			c.recordMapper.FromVendor,
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
