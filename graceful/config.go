package graceful

import "time"

type Config struct {
	ShutdownTimeout time.Duration `json:"shutdown_timeout" yaml:"shutdown_timeout" validate:"min=200ms"` //nolint:lll
	CloseTimeout    time.Duration `json:"close_timeout" yaml:"close_timeout" validate:"min=100ms"`
}

func (c Config) NewCloser(opts ...CloserOption) (Closer, error) {
	defaults := []CloserOption{
		WithCloserTimeout(c.CloseTimeout),
	}

	opts = append(defaults, opts...)

	return NewCloser(opts...)
}
