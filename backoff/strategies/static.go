package strategies

import "time"

func NewStatic(d time.Duration) Interface {
	//nolint:errcheck
	s, _ := NewExponential(
		WithMinDelay(d),
		WithMaxDelay(d),
		WithMultiplier(1),
		WithJitter(0),
	)

	return s
}
