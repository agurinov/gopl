package crontab

import "time"

type (
	Config struct {
		Jobs map[string]JobConfig `validate:"gt=0,dive,keys,required,endkeys,required"`
	}
	JobConfig struct {
		Schedule string        `validate:"cron"`
		Timeout  time.Duration `validate:"min=200ms"`
		Enabled  bool
	}
)
