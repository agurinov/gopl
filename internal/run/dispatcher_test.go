package run_test

import (
	"context"
	"io"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	irun "github.com/agurinov/gopl/internal/run"
	"github.com/agurinov/gopl/run"
	pl_testing "github.com/agurinov/gopl/testing"
)

type work struct {
	entered *atomic.Uint32
	failed  *atomic.Uint32
	exited  *atomic.Uint32
}

var (
	allKeys = []int32{5, 3, 1}
	noKeys  = []int32{}
)

func (w work) do(err error) run.Fn {
	return func(ctx context.Context) error {
		w.entered.Add(1)
		defer w.exited.Add(1)

		if err != nil {
			w.failed.Add(1)

			return err
		}

		<-ctx.Done()

		return ctx.Err()
	}
}

func newWork() work {
	return work{
		failed:  new(atomic.Uint32),
		entered: new(atomic.Uint32),
		exited:  new(atomic.Uint32),
	}
}

func TestDispatcher_Failed(t *testing.T) {
	pl_testing.Init(t)

	ctx := t.Context()

	dispatcher, err := irun.NewDispatcher[int32]()
	require.NoError(t, err)
	require.NotNil(t, dispatcher)

	var (
		w  = newWork()
		fn = w.do(io.EOF)
	)

	dispatcher.Run(ctx, fn, allKeys...)

	requireFailed(t, 3, w)
	requireExited(t, 3, w)
	require.ElementsMatch(t, dispatcher.Running(), noKeys)

	for i := range allKeys {
		iCtx, _ := dispatcher.GetContext(allKeys[i])
		require.Equal(t,
			io.EOF,
			context.Cause(iCtx),
		)
	}
}

func TestDispatcher(t *testing.T) {
	pl_testing.Init(t)

	ctx := t.Context()

	dispatcher, err := irun.NewDispatcher[int32]()
	require.NoError(t, err)
	require.NotNil(t, dispatcher)

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
			fn = w.do(nil)
		)

		dispatcher.Run(ctx, fn, allKeys...)

		require.ElementsMatch(t, dispatcher.Running(), allKeys)
		requireEntered(t, 3, w)

		if i == 0 {
			continue
		}

		previousW := iterations[i-1]
		requireExited(t, 3, previousW)
	}

	dispatcher.Stop(3)
	require.ElementsMatch(t, dispatcher.Running(), []int32{1, 5})

	dispatcher.Stop(1)
	require.ElementsMatch(t, dispatcher.Running(), []int32{5})

	dispatcher.Stop(5)
	require.ElementsMatch(t, dispatcher.Running(), noKeys)
}

func requireEntered(
	t *testing.T,
	n uint32,
	w work,
) {
	t.Helper()

	require.Eventually(
		t,
		func() bool { return w.entered.Load() == n },
		200*time.Millisecond,
		10*time.Millisecond,
		"worker never worked while expected",
	)
}

func requireExited(
	t *testing.T,
	n uint32,
	w work,
) {
	t.Helper()

	require.Eventually(
		t,
		func() bool { return w.exited.Load() == n },
		200*time.Millisecond,
		10*time.Millisecond,
		"worker can't exited; leak",
	)
}

func requireFailed(
	t *testing.T,
	n uint32,
	w work,
) {
	t.Helper()

	require.Eventually(
		t,
		func() bool { return w.failed.Load() == n },
		200*time.Millisecond,
		10*time.Millisecond,
		"worker never failed while expected",
	)
}
