package errors_test

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"syscall"
	"testing"

	"github.com/stretchr/testify/require"

	pl_errors "github.com/agurinov/gopl/errors"
	pl_testing "github.com/agurinov/gopl/testing"
)

type timeoutError struct{}

func (e timeoutError) Error() string { return "" }
func (e timeoutError) Timeout() bool { return true }

func TestIsRetryable(t *testing.T) {
	pl_testing.Init(t)

	cases := map[string]struct {
		inputErr            error
		checkWrappedGo113   bool
		expectedIsRetryable bool
		pl_testing.TestCase
	}{
		"case00: nil err": {
			inputErr:            nil,
			expectedIsRetryable: false,
		},
		"case01: io.EOF": {
			inputErr:            io.EOF,
			checkWrappedGo113:   true,
			expectedIsRetryable: true,
		},
		"case02: context.Canceled": {
			inputErr:            context.Canceled,
			checkWrappedGo113:   true,
			expectedIsRetryable: true,
		},
		"case03: context.DeadlineExceeded": {
			inputErr:            context.DeadlineExceeded,
			checkWrappedGo113:   true,
			expectedIsRetryable: true,
		},
		"case04: dns error retryable": {
			inputErr:            &net.DNSError{IsTemporary: true},
			checkWrappedGo113:   true,
			expectedIsRetryable: true,
		},
		"case04: dns error permanent": {
			inputErr:            &net.DNSError{IsTemporary: false},
			checkWrappedGo113:   true,
			expectedIsRetryable: false,
		},
		"case05: timeout error": {
			inputErr:            os.NewSyscallError("syscall_name", timeoutError{}),
			checkWrappedGo113:   true,
			expectedIsRetryable: true,
		},
		"case06: syscall.ETIMEDOUT": {
			inputErr:            syscall.ETIMEDOUT,
			checkWrappedGo113:   true,
			expectedIsRetryable: true,
		},
		"case07: syscall.EINTR": {
			inputErr:            syscall.EINTR,
			checkWrappedGo113:   true,
			expectedIsRetryable: true,
		},
		"case08: syscall.EINTR": {
			inputErr:            syscall.EADDRINUSE,
			checkWrappedGo113:   true,
			expectedIsRetryable: false,
		},
	}

	for name := range cases {
		tc := cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			require.Equal(t,
				tc.expectedIsRetryable,
				pl_errors.IsRetryable(tc.inputErr),
			)

			if !tc.checkWrappedGo113 {
				return
			}

			require.Equal(t,
				tc.expectedIsRetryable,
				pl_errors.IsRetryable(
					fmt.Errorf("%w", tc.inputErr),
				),
			)
		})
	}
}
