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
		name       string
		retries    uint32
		maxRetries uint32
		logLevel   zapcore.Level
	}
	Option = c.Option[Backoff]
)

func (b *Backoff) Wait(ctx context.Context) (Stat, error) {
	// Register new retry
	retries := atomic.AddUint32(&b.retries, 1)

	// Check limit of allowed retries
	if retries > b.maxRetries {
		return EmptyStat, RetryLimitError{
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
		return EmptyStat, ctx.Err()
	case <-time.After(delay):
		return Stat{
			Duration:   delay,
			RetryIndex: retries,
			MaxRetries: b.maxRetries,
		}, nil
	}
}

func (b *Backoff) Reset() {
	atomic.StoreUint32(&b.retries, 0)
}

func New(opts ...Option) (*Backoff, error) {
	b := Backoff{
		name:       pl_strings.UnspecifiedPlaceholder,
		maxRetries: math.MaxUint32,
		logLevel:   zapcore.InfoLevel,
	}

	obj, err := c.ConstructWithValidate(b, opts...)
	if err != nil {
		return nil, err
	}

	return &obj, nil
}
