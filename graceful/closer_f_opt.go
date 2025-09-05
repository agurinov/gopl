package graceful

import (
	"time"

	"go.uber.org/zap"
)

func WithLogger(logger *zap.Logger) CloserOption {
	return func(c *Closer) error {
		if logger == nil {
			return nil
		}

		c.logger = logger.Named("closer")

		return nil
	}
}

func WithTimeout(timeout time.Duration) CloserOption {
	return func(c *Closer) error {
		c.timeout = timeout

		return nil
	}
}
