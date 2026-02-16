package kafka

import (
	"time"

	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/agurinov/gopl/backoff"
	"github.com/agurinov/gopl/backoff/strategies"
)

type (
	ConsumerConfig struct {
		Group string
		config
		DLQ                 ProducerConfig
		MaxPollRecords      int
		MaxPollDuration     time.Duration
		idleBackoffMinDelay time.Duration
		idleBackoffMaxDelay time.Duration
	}
)

func (c ConsumerConfig) NewConsumer(
	opts ...ConsumerOption,
) (
	Consumer,
	error,
) {
	backoffOptions := []backoff.Option{
		backoff.WithName("idle"),
		backoff.WithExponentialStrategy(
			strategies.WithMinDelay(c.idleBackoffMinDelay),
			strategies.WithMaxDelay(c.idleBackoffMaxDelay),
		),
	}

	kgoConsumerOptions := []kgo.Opt{
		kgo.SeedBrokers(c.Brokers...),
		kgo.ConsumeTopics(c.Topic),
		kgo.ConsumerGroup(c.Group),
		kgo.FetchMaxWait(c.MaxPollDuration), // TODO: doesn't work
	}

	defaults := []ConsumerOption{
		WithConsumerTopic(c.Topic),
		WithBackoffOptions(backoffOptions...),
		WithConsumerClientOptions(kgoConsumerOptions...),
	}

	if c.DLQ.Topic != "" {
		// TODO: minimaze config in this case (merge base configs?)
		dlq, err := c.DLQ.NewProducer()
		if err != nil {
			return nil, err
		}

		defaults = append(defaults,
			WithDLQ(dlq),
		)
	}

	opts = append(defaults, opts...)

	return NewConsumer(opts...)
}
