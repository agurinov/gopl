package probes

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/agurinov/gopl/diag/log"
	pl_testing "github.com/agurinov/gopl/testing"
)

func TestProber_Close(t *testing.T) {
	pl_testing.Init(t)

	ctx := context.TODO()

	prober, err := New(
		WithLogger(log.NewZapTest(t)),
		WithCheckInterval(time.Second),
		WithCheckTimeout(time.Second),
		WithReadinessProbe(
			func(context.Context) error { return nil },
		),
	)
	require.NoError(t, err)
	require.NotNil(t, prober)

	prober.SetStartup(true)
	prober.runAllProbes(ctx)

	require.True(t, prober.Startup())
	require.True(t, prober.Readiness())
	require.True(t, prober.Liveness())

	prober.Close()
	prober.runAllProbes(ctx)

	require.True(t, prober.Startup())
	require.False(t, prober.Readiness())
	require.True(t, prober.Liveness())
}
