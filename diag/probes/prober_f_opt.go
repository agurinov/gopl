package probes

import (
	"errors"
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

func WithProbe(t ProbeType, probe Probe) Option {
	return func(p **Prober) error {
		pr := *p

		switch t {
		case ProbeTypeReadiness:
			pr.readinessProbes = append(pr.readinessProbes, probe)
		case ProbeTypeLiveness:
			pr.livenessProbes = append(pr.livenessProbes, probe)
		case ProbeTypeStartup:
		default:
			return errors.New("unknown probe type")
		}

		return nil
	}
}
