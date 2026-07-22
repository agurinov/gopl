package kafka

import (
	"time"

	"github.com/agurinov/gopl/backoff/strategies"
)

type (
	//nolint:lll
	ConsumerConfig struct {
		Group           string
		Topic           string   `validate:"required"`
		Brokers         []string `validate:"required,gt=0,dive,required"`
		DLQ             ProducerConfig
		MaxPollRecords  int           `json:"max_poll_records" yaml:"max_poll_records" validate:"required"`
		MaxPollDuration time.Duration `json:"max_poll_duration" yaml:"max_poll_duration"`
		Idle            strategies.ExponentialConfig
	}
)
