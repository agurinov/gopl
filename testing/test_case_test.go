package pl_testing_test

import (
	"io"
	"testing"

	pl_testing "github.com/agurinov/gopl.git/testing"
)

func TestCase_Init_CheckError(t *testing.T) {
	pl_testing.Init(t)

	cases := map[string]struct {
		err error
		pl_testing.TestCase
	}{
		"failed but skipped": {
			err: io.EOF,
			TestCase: pl_testing.TestCase{
				Skip:     true,
				MustFail: false,
			},
		},
		"failed and checked IS": {
			err: io.EOF,
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: io.EOF,
			},
		},
		"failed and checked AS": {
			err: pl_testing.ErrViolationAs,
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailAsErr: &pl_testing.TestCaseViolationError{},
			},
		},
		"success": {},
	}

	for name, tc := range cases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			tc.CheckError(t, tc.err)
		})
	}
}
