package run_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	irun "github.com/agurinov/gopl/internal/run"
	"github.com/agurinov/gopl/run"
	pl_testing "github.com/agurinov/gopl/testing"
)

type work struct {
	worked *atomic.Bool
	exited *atomic.Bool
}

func newWork() work {
	return work{
		worked: new(atomic.Bool),
		exited: new(atomic.Bool),
	}
}

func doWork(w work) run.Fn {
	return func(ctx context.Context) error {
		defer w.exited.Store(true)

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				w.worked.Store(true)
			}
		}
	}
}

func TestDispatcher(t *testing.T) {
	pl_testing.Init(t)

	ctx := t.Context()

	dispatcher, err := irun.NewDispatcher[int32]()
	require.NoError(t, err)
	require.NotNil(t, dispatcher)

	keys := []int32{5, 3, 1}

	// Run iterations with the same keys several times.
	// Expecting to have only latest goroutines set active.
	iterations := []work{
		newWork(),
		newWork(),
		newWork(),
	}

	for i := range iterations {
		var (
			w  = iterations[i]
			fn = doWork(w)
		)

		dispatcher.Run(ctx, fn, keys...)

		require.ElementsMatch(t, dispatcher.Running(), keys)
		require.Eventually(
			t,
			w.worked.Load,
			200*time.Millisecond,
			10*time.Millisecond,
			"worker can't worked",
		)

		if i == 0 {
			continue
		}

		previousW := iterations[i-1]

		require.Eventually(
			t,
			previousW.exited.Load,
			200*time.Millisecond,
			10*time.Millisecond,
			"worker can't exited; leak",
		)
	}

	dispatcher.Stop(3)
	require.ElementsMatch(t, dispatcher.Running(), []int32{1, 5})

	dispatcher.Stop(1)
	require.ElementsMatch(t, dispatcher.Running(), []int32{5})

	dispatcher.Stop(5)
	require.ElementsMatch(t, dispatcher.Running(), []int32{})
}
