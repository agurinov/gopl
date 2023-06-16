package backoff

import c "github.com/agurinov/gopl/patterns/creational"

type BackoffOption = c.Option[Backoff]

func WithStrategy(s Strategy) BackoffOption {
	return func(b *Backoff) error {
		b.strategy = s

		return nil
	}
}

func WithMaxRetries(mr uint32) BackoffOption {
	return func(b *Backoff) error {
		b.maxRetries = mr

		return nil
	}
}

func WithName(name string) BackoffOption {
	return func(b *Backoff) error {
		b.name = name

		return nil
	}
}
