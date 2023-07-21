package testing

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type TestCaseViolationError struct {
	Condition string
}

var (
	ErrViolationIs = TestCaseViolationError{Condition: "Is"}
	ErrViolationAs = TestCaseViolationError{Condition: "As"}
)

func (e TestCaseViolationError) Error() string {
	return "test case error violates " + e.Condition + " condition"
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
