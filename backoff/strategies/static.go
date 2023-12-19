package strategies

import "time"

type static time.Duration

func (s static) Duration(_ uint32) time.Duration {
	return time.Duration(s)
}

func NewStatic(d time.Duration) Interface {
	return static(d)
}
