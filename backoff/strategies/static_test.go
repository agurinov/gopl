package strategies_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/agurinov/gopl/backoff/strategies"
	pl_testing "github.com/agurinov/gopl/testing"
)

func TestStatic(t *testing.T) {
	pl_testing.Init(t)

	strategy := strategies.NewStatic(
		time.Second,
	)
	require.NotNil(t, strategy)

	cases := map[string]struct {
		pl_testing.TestCase
		expectedDuration time.Duration
		inputRetries     uint32
	}{
		"0":   {inputRetries: 0, expectedDuration: time.Second},
		"1":   {inputRetries: 1, expectedDuration: time.Second},
		"2":   {inputRetries: 2, expectedDuration: time.Second},
		"3":   {inputRetries: 3, expectedDuration: time.Second},
		"4":   {inputRetries: 4, expectedDuration: time.Second},
		"5":   {inputRetries: 5, expectedDuration: time.Second},
		"100": {inputRetries: 100, expectedDuration: time.Second},
	}

	for name := range cases {
		tc := cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			require.Equal(t,
				tc.expectedDuration,
				strategy.Duration(tc.inputRetries),
			)
		})
	}
}
