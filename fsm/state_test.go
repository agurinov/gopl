package fsm_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/agurinov/gopl/fsm"
	pl_testing "github.com/agurinov/gopl/testing"
)

func TestState_Equal(t *testing.T) {
	pl_testing.Init(t)

	type (
		args struct {
			currentState fsm.State
			otherState   fsm.State
		}
		results struct {
			isEqual bool
		}
	)

	cases := map[string]struct {
		pl_testing.TestCase
		args    args
		results results
	}{
		"case00: empty": {
			args: args{
				currentState: fsm.State{},
				otherState:   fsm.EmptyState,
			},
			results: results{
				isEqual: true,
			},
		},
		"case01: explicit equal": {
			args: args{
				currentState: fsm.State{Name: "foo"},
				otherState:   fsm.State{Name: "foo"},
			},
			results: results{
				isEqual: true,
			},
		},
		"case02: explicit different": {
			args: args{
				currentState: fsm.State{Name: "foo"},
				otherState:   fsm.State{Name: "bar"},
			},
			results: results{
				isEqual: false,
			},
		},
	}

	for name := range cases {
		tc := cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			isEqual := tc.args.currentState.Equal(tc.args.otherState)

			require.Equal(t,
				tc.results.isEqual,
				isEqual,
			)
		})
	}
}
