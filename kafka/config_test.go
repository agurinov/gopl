//go:build test_unit

package kafka_test

import (
	"testing"

	"github.com/agurinov/gopl/kafka"
	pl_testing "github.com/agurinov/gopl/testing"
)

func TestConfig_Validate(t *testing.T) {
	pl_testing.Init(t)

	cases := map[string]struct {
		inputConfig kafka.Config
		pl_testing.TestCase
	}{
		"case00: empty topic": {
			inputConfig: kafka.Config{
				BatchSize: 10,
				EventPosition: kafka.EventPosition{
					Topic: "",
				},
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: kafka.ErrInvalidConfig,
			},
		},
		"case01: zero batch size": {
			inputConfig: kafka.Config{
				BatchSize: 0,
				EventPosition: kafka.EventPosition{
					Topic: "topic_1",
				},
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: kafka.ErrInvalidConfig,
			},
		},
	}

	for name, tc := range cases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			tc.CheckError(t,
				tc.inputConfig.Validate(),
			)
		})
	}
}
