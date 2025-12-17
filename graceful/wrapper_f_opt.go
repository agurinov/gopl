package graceful

import (
	"go.uber.org/zap"
)

func WithWrapperLogger(logger *zap.Logger) WrapperOption {
	return func(c *Wrapper) error {
		if logger == nil {
			return nil
		}

		c.logger = logger.Named("graceful.wrapper")

		return nil
	}
}
