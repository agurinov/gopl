package crontab

import "time"

type (
	Config struct {
		Jobs []JobConfig `validate:"gt=0,dive"`
	}
	JobConfig struct {
		Name     string        `validate:"required"`
		Schedule string        `validate:"cron"`
		Timeout  time.Duration `validate:"min=200ms"`
	}
)
