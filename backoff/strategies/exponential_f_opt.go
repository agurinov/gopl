package strategies

import (
	"time"
)

func WithMinDelay(d time.Duration) ExponentialOption {
	return func(e *exponential) error {
		e.minDelay = d

		return nil
	}
}

func WithMaxDelay(d time.Duration) ExponentialOption {
	return func(e *exponential) error {
		e.maxDelay = d

		return nil
	}
}

func WithMultiplier(m float64) ExponentialOption {
	return func(e *exponential) error {
		e.multiplier = m

		return nil
	}
}

func WithJitter(j float64) ExponentialOption {
	return func(e *exponential) error {
		e.jitter = j

		return nil
	}
}
