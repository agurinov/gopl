package graceful_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/agurinov/gopl/graceful/internal"
	pl_testing "github.com/agurinov/gopl/testing"
)

const (
	iterationDuration        = 90 * time.Millisecond
	gracefulShutdownDuration = 100 * time.Millisecond
)

func noGraceIteration(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(iterationDuration):
		return nil
	}
}

func noGraceLoop(ctx context.Context) error {
	for {
		if err := noGraceIteration(ctx); err != nil {
			return err
		}
	}
}

func TestWrapper(t *testing.T) {
	pl_testing.Init(t)

	require.Less(t, iterationDuration, gracefulShutdownDuration)

	ctx := t.Context()

	closedCtx, cancel := context.WithCancel(ctx)
	cancel()

	t.Run("no grace", func(t *testing.T) {
		pl_testing.Init(t)

		t.Run("iteration", func(t *testing.T) {
			pl_testing.Init(t)

			require.ErrorIs(
				t,
				noGraceIteration(closedCtx),
				context.Canceled,
			)
		})

		t.Run("loop", func(t *testing.T) {
			pl_testing.Init(t)

			require.ErrorIs(
				t,
				noGraceLoop(closedCtx),
				context.Canceled,
			)
		})
	})

	t.Run("with grace for loop", func(t *testing.T) {
		pl_testing.Init(t)

		wrapper, err := internal.NewWrapper()
		require.NoError(t, err)
		require.NotNil(t, wrapper)

		t.Run("Close", func(t *testing.T) {
			pl_testing.Init(t)

			graceCtx, cancel := context.WithTimeout(ctx, gracefulShutdownDuration)
			t.Cleanup(cancel)

			closeFn := wrapper.Close(nil)
			require.NotNil(t, closeFn)

			require.NoError(t,
				closeFn(graceCtx),
			)
		})

		t.Run("Run", func(t *testing.T) {
			pl_testing.Init(t)

			t.Skip("so far we can't pass cordon to blackbox loop")

			runFn := wrapper.Run(noGraceLoop)
			require.NotNil(t, runFn)

			require.NoError(t,
				runFn(closedCtx),
			)
		})
	})

	t.Run("with grace for iteration", func(t *testing.T) {
		pl_testing.Init(t)

		wrapper, err := internal.NewWrapper()
		require.NoError(t, err)
		require.NotNil(t, wrapper)

		t.Run("Close", func(t *testing.T) {
			pl_testing.Init(t)

			graceCtx, cancel := context.WithTimeout(ctx, gracefulShutdownDuration)
			t.Cleanup(cancel)

			closeFn := wrapper.Close(nil)
			require.NotNil(t, closeFn)

			require.NoError(t,
				closeFn(graceCtx),
			)
		})

		t.Run("RunLoop", func(t *testing.T) {
			pl_testing.Init(t)

			runFn := wrapper.RunLoop(noGraceIteration)
			require.NotNil(t, runFn)

			require.NoError(t,
				runFn(closedCtx),
			)
		})
	})
}
