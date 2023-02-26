//go:build test_unit

package backoff_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/agurinov/gopl.git/backoff"
	pl_testing "github.com/agurinov/gopl.git/testing"
)

func TestStrategy_Exponential_WithoutJitter(t *testing.T) {
	pl_testing.Init(t)

	strategy, err := backoff.NewExponentialStrategy(
		backoff.WithMinDelay(1*time.Second),
		backoff.WithMaxDelay(10*time.Second),
		backoff.WithMultiplier(2.0),
		backoff.WithJitter(0.0),
	)
	require.NoError(t, err)

	cases := map[string]struct {
		inputRetries     uint32
		expectedDuration time.Duration
		pl_testing.TestCase
	}{
		"0 - MinDelay":   {inputRetries: 0, expectedDuration: 1 * time.Second},
		"1":              {inputRetries: 1, expectedDuration: 1 * time.Second},
		"2":              {inputRetries: 2, expectedDuration: 2 * time.Second},
		"3":              {inputRetries: 3, expectedDuration: 4 * time.Second},
		"4":              {inputRetries: 4, expectedDuration: 8 * time.Second},
		"5 - MaxDelay":   {inputRetries: 5, expectedDuration: 10 * time.Second},
		"100 - MaxDelay": {inputRetries: 100, expectedDuration: 10 * time.Second},
	}

	for name, tc := range cases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			require.Equal(t,
				tc.expectedDuration,
				strategy.Duration(tc.inputRetries),
			)
		})
	}
}

func TestStrategy_Exponential_WithJitter(t *testing.T) {
	pl_testing.Init(t)

	strategy, err := backoff.NewExponentialStrategy(
		backoff.WithMinDelay(1*time.Second),
		backoff.WithMaxDelay(10*time.Second),
		backoff.WithMultiplier(1.5),
		backoff.WithJitter(0.2),
	)
	require.NoError(t, err)

	cases := map[string]struct {
		inputRetries       uint32
		expectedLowerBound float64
		expectedUpperBound float64
		pl_testing.TestCase
	}{
		"0 - MinDelay": {
			inputRetries:       0,
			expectedLowerBound: float64(1) * float64(time.Second),
			expectedUpperBound: float64(1) * float64(time.Second),
		},
		"1": {
			inputRetries:       1,
			expectedLowerBound: float64(0.8) * float64(time.Second),
			expectedUpperBound: float64(1.2) * float64(time.Second),
		},
		"5 - MaxDelay": {
			inputRetries:       5,
			expectedLowerBound: float64(4.05) * float64(time.Second),
			expectedUpperBound: float64(6.075) * float64(time.Second),
		},
		"8 - MaxDelay": {
			inputRetries:       8,
			expectedLowerBound: float64(10) * float64(time.Second),
			expectedUpperBound: float64(10) * float64(time.Second),
		},
	}

	for name, tc := range cases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			duration := strategy.Duration(tc.inputRetries)
			durationFloat64 := float64(duration)

			require.GreaterOrEqual(t, durationFloat64, tc.expectedLowerBound)
			require.LessOrEqual(t, durationFloat64, tc.expectedUpperBound)
		})
	}
}
