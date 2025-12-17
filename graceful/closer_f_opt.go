package graceful

import (
	"time"

	"go.uber.org/zap"
)

func WithCloserLogger(logger *zap.Logger) CloserOption {
	return func(c *Closer) error {
		if logger == nil {
			return nil
		}

		c.logger = logger.Named("graceful.closer")

		return nil
	}
}

func WithCloserTimeout(timeout time.Duration) CloserOption {
	return func(c *Closer) error {
		c.timeout = timeout

		return nil
	}
}

func InFirstWave() AddOption {
	return func(a *addArgs) error {
		a.wave = FirstWave

		return nil
	}
}
