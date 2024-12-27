package nopanic

import (
	"go.uber.org/zap"

	"github.com/agurinov/gopl/diag/metrics"
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

func WithMetrics(enabled bool) Option {
	if !enabled {
		return nil
	}

	return func(s *Handler) error {
		s.metrics = handlerMetrics{
			panicRecovered: metrics.NewCounter(
				metrics.NopanicHandlerCounterName,
				nil,
				metrics.WithoutServicePrefix(),
			),
		}

		return nil
	}
}
