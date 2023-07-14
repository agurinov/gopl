package pl_testing

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/agurinov/gopl/env/envvars"
)

type TestCase struct {
	MustFailIsErr error
	MustFailAsErr any
	MustFail      bool

	Skip  bool
	Debug bool
	Fail  bool
}

func (tc TestCase) Init(
	t *testing.T,
	stands ...Stand,
) StandState {
	t.Helper()

	var (
		isIntegration = len(stands) > 0
		needDebug     = envvars.GDebug.Present()
		needParallel  = !needDebug
	)

	switch {
	case tc.Skip:
		t.Skip("tc skipped: explicit skip flag")
	case needDebug && !tc.Debug:
		t.Skip("tc skipped: not debuggable during " + envvars.GDebug.String())
	case testing.Short() && isIntegration:
		t.Skip("tc skipped: integration test during -testing.short mode")
	case tc.Fail:
		t.Fail()
	}

	if needParallel {
		t.Parallel()
	}

	cleanup := func() {
		// TODO(a.gurinov): deal with TestMain func
		// it doesn't work with parallel tests
		// goleak.VerifyNone(t)
		// Maybe bind it to needDebug var?
	}
	t.Cleanup(cleanup)

	creation := make(StandState, len(stands))

	for _, stand := range stands {
		created := stand.Up(t)

		if created {
			creation[stand.Name()] = struct{}{}
		}
	}

	return creation
}

func (tc TestCase) CheckError(t *testing.T, err error) {
	t.Helper()

	if !tc.MustFail {
		require.NoError(t, err)

		return
	}

	require.Error(t, err)

	if isErr := tc.MustFailIsErr; isErr != nil {
		require.ErrorIs(t, err, isErr, ErrViolationIs)
	}

	if asErr := tc.MustFailAsErr; asErr != nil {
		require.ErrorAs(t, err, asErr, ErrViolationAs)
	}

	t.Skip("tc skipped: checks after CheckError() with MustFail=true are not relevant")
}

func Init(
	t *testing.T,
	stands ...Stand,
) StandState {
	t.Helper()

	tc := TestCase{
		Debug: true,
	}

	return tc.Init(t, stands...)
}
