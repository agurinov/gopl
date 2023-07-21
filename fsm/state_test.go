package fsm_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/agurinov/gopl/fsm"
	pl_testing "github.com/agurinov/gopl/testing"
)

func TestState_Equal(t *testing.T) {
	pl_testing.Init(t)

	cases := map[string]struct {
		inputCurrentState fsm.State
		inputOtherState   fsm.State
		expectedEqual     bool
		pl_testing.TestCase
	}{
		"case00: empty": {
			inputCurrentState: fsm.State{},
			inputOtherState:   fsm.EmptyState,
			expectedEqual:     true,
		},
		"case01: explicit equal": {
			inputCurrentState: fsm.State{Name: "foo"},
			inputOtherState:   fsm.State{Name: "foo"},
			expectedEqual:     true,
		},
		"case01: explicit different": {
			inputCurrentState: fsm.State{Name: "foo"},
			inputOtherState:   fsm.State{Name: "bar"},
			expectedEqual:     false,
		},
	}

	for name, tc := range cases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			equal := tc.inputCurrentState.Equal(tc.inputOtherState)

			require.Equal(t,
				tc.expectedEqual,
				equal,
			)
		})
	}
}
