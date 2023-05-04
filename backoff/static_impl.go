package backoff

import "time"

type static time.Duration

func (s static) Duration(_ uint32) time.Duration {
	return time.Duration(s)
}

func NewStaticStrategy(d time.Duration) Strategy {
	return static(d)
}
