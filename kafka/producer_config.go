package kafka

import "github.com/twmb/franz-go/pkg/kgo"

type (
	ProducerConfig struct {
		config
	}
)

func (c ProducerConfig) NewProducer(
	opts ...ProducerOption,
) (
	Producer,
	error,
) {
	kgoProducerOptions := []kgo.Opt{
		kgo.SeedBrokers(c.Brokers...),
		kgo.ConsumeTopics(c.Topic),
	}

	defaults := []ProducerOption{
		WithProducerClient(kgoProducerOptions...),
	}

	opts = append(defaults, opts...)

	return NewProducer(opts...)
}
