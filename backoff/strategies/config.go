package strategies

import "time"

type (
	ExponentialConfig struct {
		MinDelay   time.Duration `json:"min_delay" yaml:"min_delay"`
		MaxDelay   time.Duration `json:"max_delay" yaml:"max_delay"`
		Multiplier float64
		Jitter     float64
	}
)

func (c ExponentialConfig) Options() []ExponentialOption {
	return []ExponentialOption{
		WithMinDelay(c.MinDelay),
		WithMaxDelay(c.MaxDelay),
		WithMultiplier(c.Multiplier),
		WithJitter(c.Jitter),
	}
}

func (c ExponentialConfig) New(
	opts ...ExponentialOption,
) (
	Interface,
	error,
) {
	defaults := c.Options()

	opts = append(defaults, opts...)

	return NewExponential(opts...)
}
