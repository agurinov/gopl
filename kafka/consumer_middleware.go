package kafka

import (
	"context"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/agurinov/gopl/diag/metrics"
	"github.com/agurinov/gopl/run"
	"github.com/agurinov/gopl/x"
)

var histHandlerDuration *prometheus.HistogramVec

func MetricsBatchMiddleware(
	options ...metrics.Option,
) run.Middleware[HandlerBatch] {
	options = append(
		options,
		metrics.WithoutServicePrefix(),
		metrics.WithUseExisting(),
	)
	histHandlerDuration = metrics.NewHistogram(
		metrics.KafkaConsumerHandlerBatchDurationHistogramName,
		[]string{"topic", "partition", "status", "batch"},
		options...,
	)

	return func(next HandlerBatch) HandlerBatch {
		return func(ctx context.Context, r []Record) error {
			startTime := time.Now()

			err := next(ctx, r)

			elapsedTime := time.Since(startTime)

			var (
				topic     = x.First(r).Topic
				partition = strconv.Itoa(int(x.First(r).Partition))
				status    = metrics.StatusStringFromError(err)
			)

			histHandlerDuration.WithLabelValues(
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

func MetricsMiddleware(
	options ...metrics.Option,
) run.Middleware[Handler] {
	options = append(
		options,
		metrics.WithoutServicePrefix(),
		metrics.WithUseExisting(),
	)
	histHandlerDuration = metrics.NewHistogram(
		metrics.KafkaConsumerHandlerDurationHistogramName,
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

			histHandlerDuration.WithLabelValues(
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
