package crontab

import (
	"context"
	"time"

	gocron "github.com/go-co-op/gocron/v2"
)

type Job func(context.Context) error

func taskAdapter(
	job Job,
	timeout time.Duration,
) gocron.Task {
	return gocron.NewTask(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		return job(ctx)
	})
}
