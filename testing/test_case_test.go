package testing_test

// TODO(a.gurinov): Set vars via env file for test.

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/agurinov/gopl/env/envvars"
	pl_testing "github.com/agurinov/gopl/testing"
)

func TestCase_Init_CheckError(t *testing.T) {
	pl_testing.Init(t)

	require.NoError(t,
		os.Setenv(envvars.GDebug.String(), "true"),
	)

	cases := map[string]struct {
		err error
		pl_testing.TestCase
	}{
		"failed and skipped": {
			err: io.EOF,
			TestCase: pl_testing.TestCase{
				Skip:     true,
				Debug:    true,
				MustFail: false,
			},
		},
		"failed and skipped due to not debuggable": {
			err: io.EOF,
			TestCase: pl_testing.TestCase{
				Debug:    false,
				MustFail: false,
			},
		},
		"failed and checked IS": {
			err: io.EOF,
			TestCase: pl_testing.TestCase{
				Debug:         true,
				MustFail:      true,
				MustFailIsErr: io.EOF,
			},
		},
		"failed and checked AS": {
			err: pl_testing.ErrViolationAs,
			TestCase: pl_testing.TestCase{
				Debug:         true,
				MustFail:      true,
				MustFailAsErr: &pl_testing.TestCaseViolationError{},
			},
		},
		"success": {
			TestCase: pl_testing.TestCase{
				Debug:    true,
				MustFail: false,
			},
		},
	}

	for name, tc := range cases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			tc.CheckError(t, tc.err)
		})
	}
}
