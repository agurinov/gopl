package nopanic

import (
	"go.uber.org/zap"
)

func WithLogger(logger *zap.Logger) Option {
	return func(s *Handler) error {
		if logger == nil {
			return nil
		}

		s.logger = logger.Named("nopanic.handler")

		return nil
	}
}

//revive:disable:flag-parameter
func WithMetrics(enabled bool) Option {
	if !enabled {
		return nil
	}

	return func(s *Handler) error {
		s.metrics = newHandlerMetrics()

		return nil
	}
}
