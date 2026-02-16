package kafka

import (
	"fmt"

	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
)

func WithProducerLogger(logger *zap.Logger) ProducerOption {
	return func(p *producer) error {
		if logger == nil {
			return nil
		}

		p.logger = logger.Named("kafka.producer")

		return nil
	}
}

func WithProducerTopic(topic string) ProducerOption {
	return func(p *producer) error {
		p.topic = topic

		return nil
	}
}

func WithProducerClient(opts ...kgo.Opt) ProducerOption {
	return func(p *producer) error {
		cl, err := kgo.NewClient(opts...)
		if err != nil {
			return fmt.Errorf("can't create producer kgo client: %w", err)
		}

		p.client = cl

		return nil
	}
}

//revive:disable:flag-parameter
func WithProducerMetrics(enabled bool) ProducerOption {
	if !enabled {
		return nil
	}

	return func(p *producer) error {
		p.metrics = newProducerMetrics()

		return nil
	}
}
