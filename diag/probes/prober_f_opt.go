package probes

import (
	"time"

	"go.uber.org/zap"

	"github.com/agurinov/gopl/run"
)

func WithLogger(logger *zap.Logger) Option {
	return func(p **Prober) error {
		if logger == nil {
			return nil
		}

		pr := *p
		pr.logger = logger.Named("diag.probes")

		return nil
	}
}

func WithCheckInterval(d time.Duration) Option {
	return func(p **Prober) error {
		pr := *p
		pr.checkInterval = d

		return nil
	}
}

func WithCheckTimeout(d time.Duration) Option {
	return func(p **Prober) error {
		pr := *p
		pr.checkTimeout = d

		return nil
	}
}

func WithReadinessProbe(probes ...run.Fn) Option {
	return func(p **Prober) error {
		pr := *p
		pr.WithReadinessProbe(probes...)

		return nil
	}
}

func WithLivenessProbe(probes ...run.Fn) Option {
	return func(p **Prober) error {
		pr := *p
		pr.WithLivenessProbe(probes...)

		return nil
	}
}
