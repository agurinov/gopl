package kafka

import (
	"context"
	"strconv"
	"time"

	"github.com/agurinov/gopl/diag/metrics"
	"github.com/agurinov/gopl/run"
	"github.com/agurinov/gopl/x"
)

func MetricsBatchMiddleware(
	options ...metrics.Option,
) run.Middleware[HandlerBatch] {
	options = append(
		options,
		metrics.WithoutServicePrefix(),
		metrics.WithUseExisting(),
	)
	hist := metrics.NewHistogram(
		KafkaConsumerHandlerBatchDurationHistogramName,
		[]string{"topic", "partition", "status", "batch"},
		options...,
	)

	return func(next HandlerBatch) HandlerBatch {
		return func(ctx context.Context, records []Record) error {
			startTime := time.Now()

			err := next(ctx, records)

			elapsedTime := time.Since(startTime)

			var (
				record    = x.First(records)
				topic     = record.Topic
				partition = strconv.Itoa(int(record.Partition))
				status    = metrics.StatusStringFromError(err)
				batch     = strconv.Itoa(len(records))
			)

			hist.WithLabelValues(
				topic,
				partition,
				status,
				batch,
			).Observe(
				elapsedTime.Seconds(),
			)

			return err
		}
	}
}

func MetricsMiddleware(
	options ...metrics.Option,
) run.Middleware[Handler] {
	options = append(
		options,
		metrics.WithoutServicePrefix(),
		metrics.WithUseExisting(),
	)
	hist := metrics.NewHistogram(
		KafkaConsumerHandlerDurationHistogramName,
		[]string{"topic", "partition", "status"},
		options...,
	)

	return func(next Handler) Handler {
		return func(ctx context.Context, r Record) error {
			startTime := time.Now()

			err := next(ctx, r)

			elapsedTime := time.Since(startTime)

			var (
				topic     = r.Topic
				partition = strconv.Itoa(int(r.Partition))
				status    = metrics.StatusStringFromError(err)
			)

			hist.WithLabelValues(
				topic,
				partition,
				status,
			).Observe(
				elapsedTime.Seconds(),
			)

			return err
		}
	}
}
