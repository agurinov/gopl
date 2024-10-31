package probes

import (
	"time"

	"go.uber.org/zap"
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

func WithReadinessProbe(probes ...Probe) Option {
	return func(p **Prober) error {
		pr := *p

		pr.readinessProbes = append(pr.readinessProbes, probes...)

		return nil
	}
}

func WithLivenessProbe(probes ...Probe) Option {
	return func(p **Prober) error {
		pr := *p

		pr.livenessProbes = append(pr.livenessProbes, probes...)

		return nil
	}
}
