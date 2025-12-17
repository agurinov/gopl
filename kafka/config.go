package kafka

import (
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"

	"github.com/agurinov/gopl/backoff"
	"github.com/agurinov/gopl/backoff/strategies"
)

type (
	// TODO: separate
	config struct {
		topic               string
		group               string
		brokers             []string
		maxPollRecords      int
		maxPollDuration     time.Duration
		idleBackoffMinDelay time.Duration
		idleBackoffMaxDelay time.Duration
	}
)

func (c config) kgoOptions() []kgo.Opt {
	return []kgo.Opt{
		kgo.SeedBrokers(c.brokers...),
		kgo.ConsumeTopics(c.topic),
		kgo.ConsumerGroup(c.group),
		kgo.FetchMaxWait(c.maxPollDuration), // TODO: doesn't work
	}
}

func (c config) idleBackoffOptions(logger *zap.Logger) []backoff.Option {
	return []backoff.Option{
		backoff.WithName("idle"),
		backoff.WithLogger(logger),
		backoff.WithExponentialStrategy(
			strategies.WithMinDelay(c.idleBackoffMinDelay),
			strategies.WithMaxDelay(c.idleBackoffMaxDelay),
		),
	}
}
