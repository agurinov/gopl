package backoff

import (
	"math"
	"math/rand"
	"time"

	"github.com/go-playground/validator/v10"

	c "github.com/agurinov/gopl/patterns/creational"
)

type (
	// https://en.wikipedia.org/wiki/Exponential_backoff
	// https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter
	exponential struct {
		// Global boundaries of delay duration
		minDelay time.Duration
		maxDelay time.Duration

		// multiplier is the multiplicator on each retry
		// leverages increasing of the left boundary
		// 1s, 2s, 4s, 8s, 16s (with multiplier=2.0 and minDelay=1s)
		multiplier float64

		// jitter is the randomization factor in percent J%
		// which applies boundaries on calculated backoff B
		// B -> random_from([B - J% ; B + J%])
		jitter float64
	}
	ExponentialOption = c.Option[exponential]
)

func (e exponential) Duration(retries uint32) time.Duration {
	if retries == 0 {
		return e.minDelay
	}

	// Calculate static exponential delay
	backoff := float64(e.minDelay)
	backoff *= math.Pow(e.multiplier, float64(retries-1))

	// Zero jitter means no colission avoidance
	if e.jitter != 0 {
		// https://pkg.go.dev/math/rand#Float64
		// rand.Float64() returns value from [0.0 ; 1.0)
		// So random is from [-1.0 ; 1.0)
		random := 2*rand.Float64() - 1 //nolint:gosec,gomnd

		// backoff is random from [b - delta ; b + delta)
		// where delta := b * J
		backoff *= 1 + e.jitter*random
	}

	// Check global boundaries
	min := float64(e.minDelay)
	max := float64(e.maxDelay)

	backoff = math.Max(
		min,
		math.Min(max, backoff),
	)

	return time.Duration(backoff)
}

func (e exponential) Validate() error {
	s := struct {
		MinDelay   time.Duration `validate:"min=0s"`
		MaxDelay   time.Duration `validate:"min=1s"`
		Multiplier float64       `validate:"gte=1.0"`
		Jitter     float64       `validate:"gte=0.01,lte=1.0"`
	}{
		MinDelay:   e.minDelay,
		MaxDelay:   e.maxDelay,
		Multiplier: e.multiplier,
		Jitter:     e.jitter,
	}

	if err := validator.New().Struct(s); err != nil {
		return err
	}

	return nil
}

func NewExponentialStrategy(opts ...ExponentialOption) (Strategy, error) {
	//nolint:revive,gomnd
	obj := exponential{
		minDelay:   1 * time.Second,
		maxDelay:   2 * time.Minute,
		multiplier: 1.6,
		jitter:     0.2,
	}

	return c.Construct(obj, opts...)
}
