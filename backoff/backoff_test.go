package backoff_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/agurinov/gopl/backoff"
	"github.com/agurinov/gopl/backoff/strategies"
	pl_testing "github.com/agurinov/gopl/testing"
)

func TestBackoff_Concurrency(t *testing.T) {
	pl_testing.Init(t)

	t.Run("Concurrent", func(t *testing.T) {
		var (
			maxRetries = 10
			ctx        = context.TODO()
		)

		b, err := backoff.New(
			backoff.WithName("foobar"),
			backoff.WithMaxRetries(uint32(maxRetries)),
			backoff.WithExponentialStrategy(
				strategies.WithMinDelay(1*time.Millisecond),
				strategies.WithMaxDelay(time.Second),
				strategies.WithJitter(0.0),
			),
		)
		require.NoError(t, err)
		require.NotNil(t, b)

		doValidRetries := func(ctx context.Context, b *backoff.Backoff) {
			var wg sync.WaitGroup

			wg.Add(maxRetries)

			for i := 0; i < maxRetries; i++ {
				go func() {
					defer wg.Done()

					_, err := b.Wait(ctx)
					require.NoError(t, err)
				}()
			}

			wg.Wait()
		}

		// First batch of valid retries
		doValidRetries(ctx, b)

		// Next retry in out of range
		stat, err := b.Wait(ctx)
		require.ErrorIs(t,
			err,
			backoff.RetryLimitError{
				BackoffName: "foobar",
				MaxRetries:  uint32(maxRetries),
			},
		)
		require.Equal(t, backoff.EmptyStat, stat)

		// Reset and repeat
		b.Reset()
		doValidRetries(ctx, b)
	})

	t.Run("Context", func(t *testing.T) {
		b, err := backoff.New(
			backoff.WithName("lolkek"),
			backoff.WithMaxRetries(1),
			backoff.WithExponentialStrategy(
				strategies.WithMinDelay(1*time.Hour),
				strategies.WithMaxDelay(10*time.Hour),
			),
		)
		require.NoError(t, err)
		require.NotNil(t, b)

		ctx, cancel := context.WithTimeout(
			context.Background(),
			10*time.Millisecond,
		)
		t.Cleanup(cancel)

		stat, err := b.Wait(ctx)
		require.ErrorIs(t, err, context.DeadlineExceeded)
		require.Equal(t, backoff.EmptyStat, stat)
	})
}
