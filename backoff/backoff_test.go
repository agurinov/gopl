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

func TestBackoff_Concurrency(t *testing.T) {
	pl_testing.Init(t)

	var (
		maxRetries = 10
		ctx        = context.TODO()
	)

	s, err := backoff.NewExponentialStrategy(
		backoff.WithMinDelay(1*time.Millisecond),
		backoff.WithMaxDelay(10*time.Millisecond),
		backoff.WithJitter(0.0),
	)
	require.NoError(t, err)
	require.NotNil(t, s)

	b, err := backoff.New(
		backoff.WithName("foobar"),
		backoff.WithMaxRetries(uint32(maxRetries)),
		backoff.WithStrategy(s),
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
}

func TestBackoff_Context(t *testing.T) {
	pl_testing.Init(t)

	s, err := backoff.NewExponentialStrategy(
		backoff.WithMinDelay(1*time.Hour),
		backoff.WithMaxDelay(10*time.Hour),
	)
	require.NoError(t, err)
	require.NotNil(t, s)

	b, err := backoff.New(
		backoff.WithName("lolkek"),
		backoff.WithMaxRetries(1),
		backoff.WithStrategy(s),
	)
	require.NoError(t, err)
	require.NotNil(t, b)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Millisecond,
	)
	t.Cleanup(cancel)

	stat, err := b.Wait(ctx)
	require.ErrorIs(t,
		err,
		context.DeadlineExceeded,
	)
	require.Equal(t, backoff.EmptyStat, stat)
}
