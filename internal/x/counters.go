package x

import (
	"sync"
	"sync/atomic"

	"github.com/agurinov/gopl/x"
)

type Counter[K comparable] struct {
	counts map[K]*atomic.Uint64
	mu     sync.RWMutex
}

func NewCounter[K comparable]() *Counter[K] {
	return &Counter[K]{
		counts: make(map[K]*atomic.Uint64),
	}
}

func (c *Counter[K]) Inc(key K) uint64 {
	c.mu.RLock()
	ptr, ok := c.counts[key]
	c.mu.RUnlock()

	if ok {
		return ptr.Add(1)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if ptr, ok = c.counts[key]; ok {
		return ptr.Add(1)
	}

	ptr = new(atomic.Uint64)
	c.counts[key] = ptr

	return ptr.Add(1)
}

func (c *Counter[K]) Value(key K) uint64 {
	c.mu.RLock()
	ptr, ok := c.counts[key]
	c.mu.RUnlock()

	if !ok {
		return 0
	}

	return ptr.Load()
}

func (c *Counter[K]) Values() map[K]uint64 {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return x.MapConvert(
		c.counts,
		func(k K, v *atomic.Uint64) (K, uint64) {
			return k, v.Load()
		},
	)
}
