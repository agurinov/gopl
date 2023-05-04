package pl_testing

import (
	"testing"

	"github.com/stretchr/testify/require"

	pl_bitset "github.com/agurinov/gopl/bitset"
	"github.com/agurinov/gopl/env/envvars"
)

type TestCase struct {
	MustFailIsErr error
	MustFailAsErr any
	MustFail      bool

	Skip  bool
	Debug bool
	Fail  bool

	flags pl_bitset.BitSet[TestCaseOption]
}

func (tc TestCase) Init(t *testing.T) {
	t.Helper()

	var (
		needDebug    = envvars.GDebug.Present()
		needParallel = !tc.flags.Has(TESTING_NO_PARALLEL) && !needDebug
		// needDotEnv   = !tc.flags.Has(TESTING_NO_DOTENV_FILE)
	)

	switch {
	case tc.Skip:
		t.Skip("tc skipped: explicit skip flag")
	case needDebug && !tc.Debug:
		t.Skip("tc skipped: not debuggable during " + envvars.GDebug.String())
	case tc.Fail:
		t.Fail()
	}

	cleanup := func() {
		// TODO(a.gurinov): deal with TestMain func
		// it doesn't work with parallel tests
		// goleak.VerifyNone(t)
		// Maybe bind it to needDebug var?
	}
	t.Cleanup(cleanup)

	// if needDotEnv {
	// 	require.NoError(t, dotenv.LoadOnce())
	// }

	if needParallel {
		t.Parallel()
	}
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

func Init(t *testing.T) TestCase {
	t.Helper()

	tc := TestCase{
		Debug: true,
	}
	tc.Init(t)

	return tc
}
