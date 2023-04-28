//go:build test_unit

package errors_test

import (
	"context"
	"io"
	"testing"

	pl_errors "github.com/agurinov/gopl/errors"
	pl_testing "github.com/agurinov/gopl/testing"
)

func TestOr(t *testing.T) {
	pl_testing.Init(t)

	cases := map[string]struct {
		inputErrors []error
		pl_testing.TestCase
	}{
		"case00: nil errors": {
			inputErrors: nil,
		},
		"case01: existing errors": {
			inputErrors: []error{nil, nil, io.EOF, context.Canceled},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: io.EOF,
			},
		},
	}

	for name, tc := range cases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			tc.CheckError(t, pl_errors.Or(tc.inputErrors...))
		})
	}
}
