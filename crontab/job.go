package crontab

import (
	"context"
	"time"

	gocron "github.com/go-co-op/gocron/v2"
)

type Job func(context.Context) error

func taskAdapter(
	ctx context.Context,
	job Job,
	timeout time.Duration,
) gocron.Task {
	return gocron.NewTask(
		func(ctx context.Context) error {
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			return job(ctx)
		},
		ctx,
	)
}
