package kafka

import (
	"context"
	"strconv"
	"time"

	"github.com/agurinov/gopl/diag/metrics"
	"github.com/agurinov/gopl/run"
)

const cooldownHeaderKey = "gopl_cooldown_process_after"

func CooldownMiddleware() run.Middleware[Handler] {
	hist := metrics.NewHistogram(
		KafkaConsumerCooldownDurationHistogramName,
		[]string{"topic", "partition"},
		metrics.WithoutServicePrefix(),
		metrics.WithUseExisting(),
	)

	return func(next Handler) Handler {
		return func(ctx context.Context, r Record) error {
			var cooldown time.Duration

			processAfterBytes, exists := HeaderByKey(r, cooldownHeaderKey)
			if exists {
				var processAfter time.Time
				if err := processAfter.UnmarshalBinary(processAfterBytes); err == nil {
					cooldown = time.Until(processAfter)
				}
			}

			if cooldown > 0 {
				time.Sleep(cooldown)

				var (
					topic     = r.Topic
					partition = strconv.Itoa(int(r.Partition))
				)

				hist.WithLabelValues(
					topic,
					partition,
				).Observe(
					cooldown.Seconds(),
				)
			}

			return next(ctx, r)
		}
	}
}
