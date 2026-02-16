package crontab

import (
	"context"
	"time"

	gocron "github.com/go-co-op/gocron/v2"
	"go.opentelemetry.io/otel/attribute"

	"github.com/agurinov/gopl/diag/trace"
	"github.com/agurinov/gopl/run"
)

func taskToGoCron(
	ctx context.Context,
	jobName string,
	job run.Fn,
	timeout time.Duration,
) gocron.Task {
	jobFn := func(ctx context.Context) error {
		ctx, span := trace.StartSpan(ctx, "crontab.job")
		defer span.End()

		span.SetAttributes(
			attribute.String("cronjob.name", jobName),
		)

		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		if err := job(ctx); err != nil {
			return trace.CatchError(span, err)
		}

		return nil
	}

	return gocron.NewTask(jobFn, ctx)
}
