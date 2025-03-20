package strategies_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/agurinov/gopl/backoff/strategies"
	pl_testing "github.com/agurinov/gopl/testing"
)

func TestExponential_WithoutJitter(t *testing.T) {
	pl_testing.Init(t)

	strategy, err := strategies.NewExponential(
		strategies.WithMinDelay(1*time.Second),
		strategies.WithMaxDelay(10*time.Second),
		strategies.WithMultiplier(2.0),
		strategies.WithJitter(0.0),
	)
	require.NoError(t, err)

	cases := map[string]struct {
		pl_testing.TestCase
		expectedDuration time.Duration
		inputRetries     uint32
	}{
		"0 - MinDelay":   {inputRetries: 0, expectedDuration: 1 * time.Second},
		"1":              {inputRetries: 1, expectedDuration: 1 * time.Second},
		"2":              {inputRetries: 2, expectedDuration: 2 * time.Second},
		"3":              {inputRetries: 3, expectedDuration: 4 * time.Second},
		"4":              {inputRetries: 4, expectedDuration: 8 * time.Second},
		"5 - MaxDelay":   {inputRetries: 5, expectedDuration: 10 * time.Second},
		"100 - MaxDelay": {inputRetries: 100, expectedDuration: 10 * time.Second},
	}

	for name := range cases {
		name, tc := name, cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			require.Equal(t,
				tc.expectedDuration,
				strategy.Duration(tc.inputRetries),
			)
		})
	}
}

func TestExponential_WithJitter(t *testing.T) {
	pl_testing.Init(t)

	strategy, err := strategies.NewExponential(
		strategies.WithMinDelay(1*time.Second),
		strategies.WithMaxDelay(10*time.Second),
		strategies.WithMultiplier(1.5),
		strategies.WithJitter(0.2),
	)
	require.NoError(t, err)

	cases := map[string]struct {
		pl_testing.TestCase
		expectedLowerBound float64
		expectedUpperBound float64
		inputRetries       uint32
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

	for name := range cases {
		name, tc := name, cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			duration := strategy.Duration(tc.inputRetries)
			durationFloat64 := float64(duration)

			require.GreaterOrEqual(t, durationFloat64, tc.expectedLowerBound)
			require.LessOrEqual(t, durationFloat64, tc.expectedUpperBound)
		})
	}
}
