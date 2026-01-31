package kafka

import (
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"

	"github.com/agurinov/gopl/backoff"
	"github.com/agurinov/gopl/run"
)

func WithConsumerLogger(logger *zap.Logger) ConsumerOption {
	return func(c *consumer) error {
		if logger == nil {
			return nil
		}

		c.logger = logger.Named("kafka.consumer")

		if c.recordDiscarder == nil {
			c.recordDiscarder = logRecordDiscarder{
				logger: c.logger,
			}
		}

		c.backoffOptions = append(c.backoffOptions,
			backoff.WithLogger(c.logger),
		)

		return nil
	}
}

func WithConsumerTopic(topic string) ConsumerOption {
	return func(c *consumer) error {
		c.topic = topic

		return nil
	}
}

func WithConsumerHandler(
	handler Handler,
	mws ...run.Middleware[Handler],
) ConsumerOption {
	return func(c *consumer) error {
		if len(mws) == 0 {
			c.handler = handler
		} else {
			c.handler = run.Middlewares[Handler](mws).Handler(handler)
		}

		return nil
	}
}

func WithConsumerHandlerBatch(
	handler HandlerBatch,
	mws ...run.Middleware[HandlerBatch],
) ConsumerOption {
	return func(c *consumer) error {
		if len(mws) == 0 {
			c.handlerBatch = handler
		} else {
			c.handlerBatch = run.Middlewares[HandlerBatch](mws).Handler(handler)
		}

		return nil
	}
}

func WithConsumerClientOptions(opts ...kgo.Opt) ConsumerOption {
	return func(c *consumer) error {
		c.clientOptions = append(c.clientOptions, opts...)

		return nil
	}
}

func WithBackoffOptions(opts ...backoff.Option) ConsumerOption {
	return func(c *consumer) error {
		c.backoffOptions = append(c.backoffOptions, opts...)

		return nil
	}
}

func WithMaxPollDuration(maxPollDuration time.Duration) ConsumerOption {
	return func(c *consumer) error {
		c.maxPollDuration = maxPollDuration

		return nil
	}
}

func WithMaxPollRecords(maxPollRecords int) ConsumerOption {
	return func(c *consumer) error {
		c.maxPollRecords = maxPollRecords

		return nil
	}
}

//revive:disable:flag-parameter
func WithConsumerMetrics(enabled bool) ConsumerOption {
	if !enabled {
		return nil
	}

	return func(c *consumer) error {
		c.metrics = newConsumerMetrics()

		return nil
	}
}

func WithDLQ(dlq Producer) ConsumerOption {
	return func(c *consumer) error {
		c.recordDiscarder = dlqRecordDiscarder{
			dlq: dlq,
		}

		return nil
	}
}
