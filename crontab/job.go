package crontab

import (
	"context"
	"time"

	gocron "github.com/go-co-op/gocron/v2"
	"go.opentelemetry.io/otel/attribute"

	"github.com/agurinov/gopl/diag/trace"
)

type Job func(context.Context) error

func taskAdapter(
	ctx context.Context,
	jobName string,
	job Job,
	timeout time.Duration,
) gocron.Task {
	jobF := func(ctx context.Context) error {
		ctx, span := trace.StartSpan(ctx, "crontab.job")
		defer span.End()

		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		span.SetAttributes(
			attribute.String("cronjob.name", jobName),
		)

		if err := job(ctx); err != nil {
			return trace.CatchError(span, err)
		}

		return nil
	}

	return gocron.NewTask(jobF, ctx)
}
