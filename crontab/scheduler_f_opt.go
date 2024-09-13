package crontab

import (
	"context"
	"errors"
	"time"

	gocron "github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func WithLogger(logger *zap.Logger) SchedulerOption {
	return func(_ context.Context, s *Scheduler) error {
		if logger == nil {
			return nil
		}

		s.logger = logger.Named("crontab.scheduler")

		return nil
	}
}

func WithJobRegistry(jobs map[string]Job) SchedulerOption {
	return func(_ context.Context, s *Scheduler) error {
		s.jobs = jobs

		return nil
	}
}

func WithShutdownTimeout(t time.Duration) SchedulerOption {
	return func(_ context.Context, s *Scheduler) error {
		s.shutdownTimeout = t

		return nil
	}
}

func withDefaultScheduler() SchedulerOption {
	return func(_ context.Context, s *Scheduler) error {
		if s.logger == nil {
			return errors.New("crontab: scheduler requires logger")
		}

		scheduler, err := gocron.NewScheduler(
			gocron.WithLogger(loggerAdapter(s.logger)),
			gocron.WithStopTimeout(s.shutdownTimeout),
			gocron.WithGlobalJobOptions(
				gocron.WithSingletonMode(gocron.LimitModeReschedule),
				gocron.WithEventListeners(
					gocron.BeforeJobRuns(func(jobUUID uuid.UUID, jobName string) {
						s.logger.Info(
							"starting job",
							zap.String("job_name", jobName),
							zap.String("job_uuid", jobUUID.String()),
						)
					}),
					gocron.AfterJobRuns(func(jobUUID uuid.UUID, jobName string) {
						s.logger.Info(
							"job finished",
							zap.String("job_name", jobName),
							zap.String("job_uuid", jobUUID.String()),
						)
					}),
					gocron.AfterJobRunsWithError(func(jobUUID uuid.UUID, jobName string, err error) {
						s.logger.Error(
							"can't finish job: error",
							zap.String("job_name", jobName),
							zap.String("job_uuid", jobUUID.String()),
							zap.Error(err),
						)
					}),
					gocron.AfterJobRunsWithPanic(func(jobUUID uuid.UUID, jobName string, recoverData any) {
						s.logger.Error(
							"can't finish job: panic",
							zap.String("job_name", jobName),
							zap.String("job_uuid", jobUUID.String()),
							zap.Any("panic", recoverData),
							zap.Stack("stack"),
						)
					}),
				),
			),
		)
		if err != nil {
			return err
		}

		s.scheduler = scheduler

		return nil
	}
}

func WithJob(cfg JobConfig) SchedulerOption {
	return func(ctx context.Context, s *Scheduler) error {
		if len(s.jobs) == 0 {
			return errors.New("crontab: no job registry provided")
		}

		job, exists := s.jobs[cfg.Name]
		if !exists {
			return errors.New("crontab: no job found in registry")
		}

		if s.scheduler == nil {
			if err := withDefaultScheduler()(ctx, s); err != nil {
				return err
			}
		}

		if _, jobErr := s.scheduler.NewJob(
			gocron.CronJob(cfg.Schedule, false),
			taskAdapter(ctx, job, cfg.Timeout),
			gocron.WithName(cfg.Name),
			gocron.WithIdentifier(
				uuid.NewMD5(uuid.Nil, []byte(cfg.Name)),
			),
		); jobErr != nil {
			return jobErr
		}

		return nil
	}
}