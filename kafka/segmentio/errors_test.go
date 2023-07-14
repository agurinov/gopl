//go:build test_unit

package segmentio_test

import (
	"errors"
	"io"
	"testing"

	segmentio "github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"

	pl_errors "github.com/agurinov/gopl/errors"
	pl_testing "github.com/agurinov/gopl/testing"
)

func TestErrors(t *testing.T) {
	pl_testing.Init(t)

	cases := map[string]struct {
		inputError          error
		expectedIsRetryable bool
		pl_testing.TestCase
	}{
		"case00: segmentio.ErrGenerationEnded": {
			inputError:          segmentio.ErrGenerationEnded,
			expectedIsRetryable: false,
		},
		"case01: segmentio.ErrGroupClosed": {
			inputError:          segmentio.ErrGroupClosed,
			expectedIsRetryable: false,
		},
		"case02: segmentio.RequestTimedOut": {
			inputError:          segmentio.RequestTimedOut,
			expectedIsRetryable: true,
		},
		"case03: segmentio.LeaderNotAvailable": {
			inputError:          segmentio.LeaderNotAvailable,
			expectedIsRetryable: true,
		},
		"case04: segmentio.WriteErrors": {
			inputError: segmentio.WriteErrors([]error{
				segmentio.RequestTimedOut,
				segmentio.LeaderNotAvailable,
				io.EOF,
				errors.New("foobar"),
			}),
			expectedIsRetryable: false,
		},
		"case05: segmentio.MessageTooLargeError": {
			inputError:          segmentio.MessageTooLargeError{},
			expectedIsRetryable: false,
		},
	}

	for name, tc := range cases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			require.Equal(t,
				tc.expectedIsRetryable,
				pl_errors.IsRetryable(tc.inputError),
			)
		})
	}
}
