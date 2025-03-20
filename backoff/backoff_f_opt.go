package backoff

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/agurinov/gopl/backoff/strategies"
)

func WithLogger(logger *zap.Logger) Option {
	return func(b *Backoff) error {
		if logger == nil {
			return nil
		}

		b.logger = logger.Named("backoff")

		return nil
	}
}

func WithLogLevel(logLevel zapcore.Level) Option {
	return func(b *Backoff) error {
		b.logLevel = logLevel

		return nil
	}
}

func WithExponentialStrategy(opts ...strategies.ExponentialOption) Option {
	return func(b *Backoff) error {
		exponentialStrategy, err := strategies.NewExponential(opts...)
		if err != nil {
			return err
		}

		b.strategy = exponentialStrategy

		return nil
	}
}

func WithStrategy(s strategies.Interface) Option {
	return func(b *Backoff) error {
		b.strategy = s

		return nil
	}
}

func WithMaxRetries(mr uint32) Option {
	return func(b *Backoff) error {
		b.maxRetries = mr

		return nil
	}
}

func WithName(name string) Option {
	return func(b *Backoff) error {
		b.name = name

		return nil
	}
}
