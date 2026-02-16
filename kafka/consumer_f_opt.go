package kafka

import (
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"

	"github.com/agurinov/gopl/backoff"
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
	mws ...Middleware,
) ConsumerOption {
	return func(c *consumer) error {
		if len(mws) == 0 {
			c.handler = handler
		} else {
			c.handler = Middlewares(mws).Handler(handler)
		}

		return nil
	}
}

func WithConsumerHandlerBatch(handlerBatch HandlerBatch) ConsumerOption {
	return func(c *consumer) error {
		c.handlerBatch = handlerBatch

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
		c.backoffFabric = opts

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
