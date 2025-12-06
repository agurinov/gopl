package kafka

import (
	"fmt"

	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
)

func WithConsumerLogger(logger *zap.Logger) ConsumerOption {
	return func(c *consumer) error {
		if logger == nil {
			return nil
		}

		c.logger = logger.Named("kafka.consumer")

		return nil
	}
}

func WithConsumerTopic(topic string) ConsumerOption {
	return func(c *consumer) error {
		c.config.topic = topic

		return nil
	}
}

func WithConsumerConnect() ConsumerOption {
	return func(c *consumer) error {
		copts := c.kgoOptions()

		cl, err := kgo.NewClient(copts...)
		if err != nil {
			return fmt.Errorf("can't create kgo client: %w", err)
		}

		c.client = cl

		return nil
	}
}
