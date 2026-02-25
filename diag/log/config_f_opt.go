package log

import "go.uber.org/zap"

func WithFormat(format string) Option {
	return func(c *Config) error {
		c.Format = format

		return nil
	}
}

func WithLevel(level string) Option {
	return func(c *Config) error {
		c.Level = level

		return nil
	}
}

func WithCaller(enableCaller bool) Option {
	return func(c *Config) error {
		c.EnableCaller = enableCaller

		return nil
	}
}

func WithTraceback(enableTraceback bool) Option {
	return func(c *Config) error {
		c.EnableTraceback = enableTraceback

		return nil
	}
}

func WithFields(zapFields ...zap.Field) Option {
	return func(c *Config) error {
		c.zapFields = zapFields

		return nil
	}
}
