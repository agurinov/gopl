//go:build test_unit

package backoff_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/agurinov/gopl.git/backoff"
	pl_testing "github.com/agurinov/gopl.git/testing"
)

func TestBackoff_Concurrent(t *testing.T) {
	pl_testing.Init(t)

	var (
		maxRetries = 10
		ctx        = context.TODO()
	)

	s, err := backoff.NewExponentialStrategy(
		backoff.WithMinDelay(1*time.Millisecond),
		backoff.WithMaxDelay(10*time.Millisecond),
	)
	require.NoError(t, err)

	b, err := backoff.New(
		backoff.WithName("foobar"),
		backoff.WithMaxRetries(uint32(maxRetries)),
		backoff.WithStrategy(s),
	)
	require.NoError(t, err)

	doValidRetries := func(ctx context.Context, b *backoff.Backoff) {
		var wg sync.WaitGroup

		wg.Add(maxRetries)

		for i := 0; i < maxRetries; i++ {
			go func() {
				defer wg.Done()

				require.NoError(t,
					b.Wait(ctx),
				)
			}()
		}

		wg.Wait()
	}

	// First batch of valid retries
	doValidRetries(ctx, b)

	// Next retry in out of range
	require.ErrorIs(t,
		b.Wait(ctx),
		backoff.RetryLimitError{
			BackoffName: "foobar",
			MaxRetries:  uint32(maxRetries),
		},
	)

	// Reset and repeat
	b.Reset()
	doValidRetries(ctx, b)
}

func TestBackoff_Context(t *testing.T) {
	pl_testing.Init(t)

	s, err := backoff.NewExponentialStrategy(
		backoff.WithMinDelay(1*time.Hour),
		backoff.WithMaxDelay(10*time.Hour),
	)
	require.NoError(t, err)

	b, err := backoff.New(
		backoff.WithName("lolkek"),
		backoff.WithMaxRetries(1),
		backoff.WithStrategy(s),
	)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Millisecond,
	)
	t.Cleanup(cancel)

	require.ErrorIs(t,
		b.Wait(ctx),
		context.DeadlineExceeded,
	)
}
