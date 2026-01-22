package kafka

import (
	"fmt"

	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
)

func WithConsumerLogger[T any](logger *zap.Logger) ConsumerOption[T] {
	return func(c *consumer[T]) error {
		if logger == nil {
			return nil
		}

		c.logger = logger.Named("kafka.consumer")

		return nil
	}
}

func WithConsumerTopic[T any](topic string) ConsumerOption[T] {
	return func(c *consumer[T]) error {
		c.config.topic = topic

		return nil
	}
}

func WithConsumerGroup[T any](group string) ConsumerOption[T] {
	return func(c *consumer[T]) error {
		c.config.group = group

		return nil
	}
}

func WithConsumerBrokers[T any](brokers []string) ConsumerOption[T] {
	return func(c *consumer[T]) error {
		c.config.brokers = brokers

		return nil
	}
}

func WithConsumerHandler[T any](handler Handler[T]) ConsumerOption[T] {
	return func(c *consumer[T]) error {
		c.handler = handler

		return nil
	}
}

func WithConsumerHandlerBatch[T any](handlerBatch HandlerBatch[T]) ConsumerOption[T] {
	return func(c *consumer[T]) error {
		c.handlerBatch = handlerBatch

		return nil
	}
}

func WithConsumerConnect[T any]() ConsumerOption[T] {
	return func(c *consumer[T]) error {
		copts := c.kgoOptions()

		cl, err := kgo.NewClient(copts...)
		if err != nil {
			return fmt.Errorf("can't create kgo client: %w", err)
		}

		c.client = cl

		return nil
	}
}
