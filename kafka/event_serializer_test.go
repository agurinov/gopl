//go:build test_unit

package kafka_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/agurinov/gopl/kafka"
	pl_testing "github.com/agurinov/gopl/testing"
)

func TestEventSerializer_JSON(t *testing.T) {
	pl_testing.Init(t)

	cases := map[string]struct {
		inputData     []byte
		expectedEvent Contract
		pl_testing.TestCase
	}{
		"case00: syntax": {
			inputData: brokenEvent1Json,
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailAsErr: NewSerializationError(),
			},
		},
		"case01: success": {
			inputData:     event1Json,
			expectedEvent: event1,
		},
	}

	for name, tc := range cases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			e, err := kafka.JsonEventSerializer[Contract](tc.inputData)

			tc.CheckError(t, err)
			require.Equal(t, tc.expectedEvent, e)
		})
	}
}

func TestEventSerializer_Proto(t *testing.T) {
	pl_testing.Init(t)

	cases := map[string]struct {
		inputData     []byte
		expectedEvent Contract
		pl_testing.TestCase
	}{
		"case00: syntax err": {
			inputData: brokenEvent1Json,
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				Skip:          true,
				MustFailAsErr: NewSerializationError(),
			},
		},
		"case01: success": {
			inputData:     event1Json,
			expectedEvent: event1,
			TestCase: pl_testing.TestCase{
				Skip: true,
			},
		},
	}

	for name, tc := range cases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			e, err := kafka.ProtoEventSerializer[Contract](tc.inputData)

			tc.CheckError(t, err)
			require.Equal(t, tc.expectedEvent, e)
		})
	}
}
