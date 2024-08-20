package crontab

import (
	"errors"
	"time"

	gocron "github.com/go-co-op/gocron/v2"
	"go.uber.org/zap"
)

func WithLogger(logger *zap.Logger) SchedulerOption {
	return func(s *Scheduler) error {
		if logger == nil {
			return nil
		}

		s.logger = logger.Named("crontab.scheduler")

		return nil
	}
}

func WithJobRegistry(jobs map[string]Job) SchedulerOption {
	return func(s *Scheduler) error {
		s.jobs = jobs

		return nil
	}
}

func WithShutdownTimeout(t time.Duration) SchedulerOption {
	return func(s *Scheduler) error {
		s.shutdownTimeout = t

		return nil
	}
}

func WithJob(cfg JobConfig) SchedulerOption {
	return func(s *Scheduler) error {
		if len(s.jobs) == 0 {
			return errors.New("crontab: no job registry provided")
		}

		job, exists := s.jobs[cfg.Name]
		if !exists {
			return errors.New("crontab: no job found in registry")
		}

		if s.scheduler == nil {
			scheduler, err := gocron.NewScheduler(
				gocron.WithLogger(loggerAdapter(s.logger)),
				gocron.WithStopTimeout(s.shutdownTimeout),
			)
			if err != nil {
				return err
			}

			s.scheduler = scheduler
		}

		if _, jobErr := s.scheduler.NewJob(
			gocron.CronJob(cfg.Schedule, false),
			taskAdapter(job, cfg.Timeout),
			gocron.WithName(cfg.Name),
			gocron.WithSingletonMode(gocron.LimitModeReschedule),
		); jobErr != nil {
			return jobErr
		}

		return nil
	}
}
