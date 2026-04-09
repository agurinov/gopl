package crontab

import (
	"context"
	"time"

	"github.com/agurinov/gopl/x"
)

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

func (c Config) New(
	ctx context.Context,
	opts ...SchedulerOption,
) (
	Scheduler,
	error,
) {
	defaults := x.MapToSlice(c.Jobs, WithJob)

	opts = append(defaults, opts...)

	return New(ctx, opts...)
}
