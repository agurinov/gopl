package backoff

import (
	"context"
	"math"
	"sync/atomic"
	"time"

	c "github.com/agurinov/gopl.git/pattern/creational"
	pl_strings "github.com/agurinov/gopl.git/strings"
)

// TODO(a.gurinov): tests
type Backoff struct {
	strategy   Strategy
	name       string
	retries    uint32
	maxRetries uint32 `validate:"min=1"`
}

func (b *Backoff) Wait(ctx context.Context) error {
	// Register new retry
	retries := atomic.AddUint32(&b.retries, 1)

	// Check limit of allowed retries
	if retries > b.maxRetries {
		return RetryLimitError{
			BackoffName: b.name,
			MaxRetries:  b.maxRetries,
		}
	}

	delay := b.strategy.Duration(retries)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(delay):
		return nil
	}
}

func (b *Backoff) Reset() {
	atomic.StoreUint32(&b.retries, 0)
}

func New(opts ...BackoffOption) (*Backoff, error) {
	exponentialStrategy, err := NewExponentialStrategy()
	if err != nil {
		return nil, err
	}

	obj := Backoff{
		name:       pl_strings.UnspecifiedPlaceholder,
		maxRetries: math.MaxUint32,
		strategy:   exponentialStrategy,
	}

	obj, err = c.ConstructObject(obj, opts...)
	if err != nil {
		return nil, err
	}

	return &obj, nil
}
