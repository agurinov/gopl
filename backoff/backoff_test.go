package backoff_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/agurinov/gopl/backoff"
	pl_testing "github.com/agurinov/gopl/testing"
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

func TestBackoff_Validate(t *testing.T) {
	pl_testing.Init(t)

	cases := map[string]struct {
		inputOptions []backoff.Option
		pl_testing.TestCase
	}{
		"case00: success": {
			inputOptions: []backoff.Option{
				backoff.WithMaxRetries(5),
			},
		},
	}

	for name := range cases {
		name, tc := name, cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			b, err := backoff.New(tc.inputOptions...)
			require.NoError(t, err)

			tc.CheckError(t,
				b.Validate(),
			)
		})
	}
}
