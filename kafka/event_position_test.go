//go:build test_unit

package kafka_test

import (
	"testing"

	"github.com/agurinov/gopl/kafka"
	pl_testing "github.com/agurinov/gopl/testing"
)

func TestEventPosition_ValidateWith(t *testing.T) {
	pl_testing.Init(t)

	configEventPosition := kafka.EventPosition{
		Topic:     "topic_1",
		Partition: 8,
		Offset:    143298,
	}

	cases := map[string]struct {
		inputConfiguredEventPosition kafka.EventPosition
		inputEventPosition           kafka.EventPosition
		pl_testing.TestCase
	}{
		"case00: unexpected topic": {
			inputConfiguredEventPosition: configEventPosition,
			inputEventPosition: kafka.EventPosition{
				Topic:     "topic_2",
				Partition: 8,
				Offset:    143298,
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: kafka.ErrInvalidEventPosition,
			},
		},
		"case01: known but unexpected partition": {
			inputConfiguredEventPosition: configEventPosition,
			inputEventPosition: kafka.EventPosition{
				Topic:     "topic_1",
				Partition: 4,
				Offset:    143298,
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: kafka.ErrInvalidEventPosition,
			},
		},
		"case02: unknown partition always ok": {
			inputConfiguredEventPosition: kafka.EventPosition{
				Topic:     "topic_1",
				Partition: kafka.UknownPartition,
			},
			inputEventPosition: kafka.EventPosition{
				Topic:     "topic_1",
				Partition: 4,
				Offset:    143298,
			},
		},
		"case03: unexpected offset": {
			inputConfiguredEventPosition: configEventPosition,
			inputEventPosition: kafka.EventPosition{
				Topic:     "topic_1",
				Partition: 8,
				Offset:    143297,
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: kafka.ErrInvalidEventPosition,
			},
		},
	}

	for name, tc := range cases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			tc.CheckError(t, tc.inputEventPosition.ValidateWith(
				tc.inputConfiguredEventPosition,
			))
		})
	}
}
