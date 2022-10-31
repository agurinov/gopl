package pl_testing

import (
	"testing"

	"github.com/stretchr/testify/require"
	_ "go.uber.org/goleak"

	pl_bitset "github.com/agurinov/gopl.git/bitset"
)

type TestCase struct {
	MustFailIsErr error
	MustFailAsErr error
	MustFail      bool

	flags pl_bitset.BitSet[TestCaseOption]
}

func (tc TestCase) Init(t *testing.T) {
	t.Helper()

	cleanup := func() {
		// TODO(a.gurinov): deal with TestMain func
		// it doesn't work with parallel tests
		// goleak.VerifyNone(t)
	}
	t.Cleanup(cleanup)

	//nolint:gofumpt
	var (
		// needDotEnv   = !tc.flags.Has(TESTING_NO_DOTENV_FILE)
		needParallel = !tc.flags.Has(TESTING_NO_PARALLEL)
	)

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

	if isErr := tc.MustFailIsErr; isErr != nil {
		require.ErrorIs(t, err, isErr, ErrViolationIs)
	}

	if asErr := tc.MustFailAsErr; asErr != nil {
		require.ErrorAs(t, err, asErr, ErrViolationAs)
	}
}

func Init(t *testing.T) {
	t.Helper()

	tc := TestCase{}
	tc.Init(t)
}
