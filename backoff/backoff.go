package backoff

import (
	"context"
	"math"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/agurinov/gopl/backoff/strategies"
	c "github.com/agurinov/gopl/patterns/creational"
	pl_strings "github.com/agurinov/gopl/strings"
)

type (
	Backoff struct {
		strategy   strategies.Interface
		logger     *zap.Logger
		retries    *atomic.Uint32
		name       string
		maxRetries uint32
		logLevel   zapcore.Level
	}
	Option = c.Option[Backoff]
)

func (b Backoff) Wait(ctx context.Context) (Stat, error) {
	// Register new retry
	retries := b.retries.Add(1)

	var (
		isUnlimited     = b.maxRetries == math.MaxUint32
		isLimitExceeded = retries > b.maxRetries
	)

	if !isUnlimited && isLimitExceeded {
		return Stat{}, RetryLimitError{
			BackoffName: b.name,
			MaxRetries:  b.maxRetries,
		}
	}

	delay := b.strategy.Duration(retries)

	b.logger.Named(b.name).Log(
		b.logLevel,
		"backoff occurred",
		zap.Uint32("retries", retries),
		zap.Stringer("delay", delay),
	)

	select {
	case <-ctx.Done():
		return Stat{}, ctx.Err()
	case <-time.After(delay):
		stat := Stat{
			Duration:   delay,
			RetryIndex: retries,
			MaxRetries: b.maxRetries,
		}

		return stat, nil
	}
}

func (b Backoff) Reset() {
	b.retries.Store(0)
}

func New(opts ...Option) (Backoff, error) {
	b := Backoff{
		name:       pl_strings.UnspecifiedPlaceholder,
		maxRetries: math.MaxUint32,
		logLevel:   zapcore.InfoLevel,
		retries:    new(atomic.Uint32),
	}

	return c.ConstructWithValidate(b, opts...)
}
