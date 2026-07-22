package kafka

import (
	"time"
)

type (
	ProducerConfig struct {
		Topic    string        `validate:"required"`
		Brokers  []string      `validate:"required,gt=0,dive,required"`
		Cooldown time.Duration `validate:"gte=0"`
	}
)
