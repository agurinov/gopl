package crontab

import (
	"context"
	"fmt"
	"time"

	gocron "github.com/go-co-op/gocron/v2"
	"go.uber.org/zap"

	c "github.com/agurinov/gopl/patterns/creational"
)

type (
	Scheduler struct {
		scheduler       gocron.Scheduler
		logger          *zap.Logger
		jobs            map[string]Job
		shutdownTimeout time.Duration
	}
	SchedulerOption c.OptionWithContext[Scheduler]
)

var New = c.NewWithContextValidate[Scheduler, SchedulerOption]

func (s Scheduler) Run(_ context.Context) error {
	s.logger.Info("starting crontab")

	s.scheduler.Start()

	return nil
}

func (s Scheduler) WaitForShutdown(ctx context.Context) error {
	<-ctx.Done()

	s.logger.Debug("shutting down crontab")

	if err := s.scheduler.Shutdown(); err != nil {
		return fmt.Errorf("crontab: can't shutdown: %w", err)
	}

	return nil
}
