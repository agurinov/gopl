package backoff

import "github.com/agurinov/gopl/backoff/strategies"

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
