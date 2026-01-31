package kafka

import (
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap/zapcore"

	"github.com/agurinov/gopl/backoff"
	"github.com/agurinov/gopl/backoff/strategies"
)

type (
	ConsumerConfig struct {
		Group           string
		Topic           string
		Brokers         []string
		DLQ             ProducerConfig
		MaxPollRecords  int           `json:"max_poll_records" yaml:"max_poll_records" validate:"required"`
		MaxPollDuration time.Duration `json:"max_poll_duration" yaml:"max_poll_duration"`
		Idle            strategies.ExponentialConfig
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
		backoff.WithLogLevel(zapcore.DebugLevel),
		backoff.WithUnlimitedRetries(),
		backoff.WithExponentialStrategy(),
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
		WithMaxPollDuration(c.MaxPollDuration),
		WithMaxPollRecords(c.MaxPollRecords),
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
