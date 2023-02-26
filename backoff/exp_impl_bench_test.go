package backoff_test

import (
	"math"
	"testing"
	"time"
)

func BenchmarkExponential(b *testing.B) {
	var (
		multiplier     = 2.0
		initialBackoff = float64(time.Second)
	)

	var (
		loopLogic = func(retries uint32) float64 {
			b := initialBackoff

			for i := 0; i < int(retries); i++ {
				b *= multiplier
			}

			return b
		}

		// Much faster than loopLogic
		mathLogic = func(retries uint32) float64 {
			b := initialBackoff

			return b * math.Pow(multiplier, float64(retries))
		}
	)

	cases := map[string]struct {
		inputRetries uint32
	}{
		"retries=10":   {inputRetries: 10},
		"retries=100":  {inputRetries: 100},
		"retries=1000": {inputRetries: 1000},
	}

	for name, bc := range cases {
		b.Run(name, func(b *testing.B) {
			b.Run("loop logic", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					loopLogic(bc.inputRetries)
				}
			})

			b.Run("math logic", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					mathLogic(bc.inputRetries)
				}
			})
		})
	}
}
