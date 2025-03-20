package strategies

import (
	"math/rand"
	"time"
)

type random struct {
	maxDelay time.Duration
}

func (r random) Duration(uint32) time.Duration {
	random := rand.Int63n(int64(r.maxDelay)) //nolint:gosec

	return time.Duration(random)
}

func NewRandom(maxDelay time.Duration) Interface {
	return random{
		maxDelay: maxDelay,
	}
}
