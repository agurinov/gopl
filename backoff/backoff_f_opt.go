package backoff

func WithStrategy(s Strategy) Option {
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
