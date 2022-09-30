package pl_testing

import (
	"testing"

	"github.com/stretchr/testify/require"

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

	var (
		needDotEnv   = !tc.flags.Has(TESTING_NO_DOTENV_FILE)
		needParallel = !tc.flags.Has(TESTING_NO_PARALLEL)
	)

	if needDotEnv {
		// require.NoError(t, dotenv.LoadOnce())
	}

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
