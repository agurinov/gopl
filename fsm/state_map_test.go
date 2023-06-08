package fsm_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/agurinov/gopl/fsm"
	pl_testing "github.com/agurinov/gopl/testing"
)

func TestStateMap_New(t *testing.T) {
	pl_testing.Init(t)

	cases := map[string]struct {
		inputStates []fsm.State
		pl_testing.TestCase
	}{
		"case00: no states at all": {
			inputStates: nil,
		},
		"case01: duplicate states": {
			inputStates: []fsm.State{
				{Name: "foo"},
				{Name: "bar"},
				{Name: "foo"},
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: fsm.StateAlreadyPresentError{StateName: "foo"},
			},
		},
		"case02: duplicate initial": {
			inputStates: []fsm.State{
				{Name: "foo", Initial: true},
				{Name: "bar", Initial: true},
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: fsm.StateAlreadyPresentError{StateName: "foo"},
			},
		},
		"case03: duplicate final": {
			inputStates: []fsm.State{
				{Name: "foo", Final: true},
				{Name: "bar", Final: true},
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: fsm.StateAlreadyPresentError{StateName: "foo"},
			},
		},
		"case04: duplicate broken": {
			inputStates: []fsm.State{
				{Name: "foo", Broken: true},
				{Name: "bar", Broken: true},
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: fsm.StateAlreadyPresentError{StateName: "foo"},
			},
		},
		"case05: success": {
			inputStates: []fsm.State{
				{Name: "foo", Initial: true},
				{Name: "bar", Broken: true},
				{Name: "baz", Final: true},
				{Name: "lol"},
			},
		},
	}

	for name, tc := range cases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			_, err := fsm.NewStateMap(tc.inputStates...)
			tc.CheckError(t, err)
		})
	}
}

func TestStateMap_Validate(t *testing.T) {
	pl_testing.Init(t)

	cases := map[string]struct {
		inputStates          []fsm.State
		expectedInitialState fsm.State
		expectedFinalState   fsm.State
		expectedBrokenState  fsm.State
		pl_testing.TestCase
	}{
		"case00: without initial": {
			inputStates: []fsm.State{
				{Name: "foo"},
				{Name: "bar", Final: true},
				{Name: "baz", Broken: true},
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: fsm.StateNotPresentError{StateName: "initial"},
			},
		},
		"case01: without final state": {
			inputStates: []fsm.State{
				{Name: "foo", Initial: true},
				{Name: "bar"},
				{Name: "baz", Broken: true},
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: fsm.StateNotPresentError{StateName: "final"},
			},
		},
		"case02: without broken state": {
			inputStates: []fsm.State{
				{Name: "foo", Initial: true},
				{Name: "bar", Final: true},
				{Name: "baz"},
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: fsm.StateNotPresentError{StateName: "broken"},
			},
		},
		"case03: initial state without next possible": {
			inputStates: []fsm.State{
				{Name: "foo", Initial: true},
				{Name: "bar", Final: true},
				{Name: "baz", Broken: true},
			},
			TestCase: pl_testing.TestCase{
				MustFail: true,
				MustFailIsErr: fsm.UnexpectedPossibleStatesError{
					StateName:     "foo",
					MustBePresent: true,
				},
			},
		},
		"case04: final state have next possible": {
			inputStates: []fsm.State{
				{
					Name: "foo", Initial: true,
					PossibleStates: fsm.MustNewStateMap(
						fsm.State{Name: "lol"},
					),
				},
				{
					Name: "bar", Final: true,
					PossibleStates: fsm.MustNewStateMap(
						fsm.State{Name: "lol"},
					),
				},
				{Name: "baz", Broken: true},
			},
			TestCase: pl_testing.TestCase{
				MustFail: true,
				MustFailIsErr: fsm.UnexpectedPossibleStatesError{
					StateName:          "bar",
					MustBePresent:      false,
					PossibleStateNames: "lol",
				},
			},
		},
		"case05: broken state have next possible": {
			inputStates: []fsm.State{
				{
					Name: "foo", Initial: true,
					PossibleStates: fsm.MustNewStateMap(
						fsm.State{Name: "lol"},
					),
				},
				{Name: "bar", Final: true},
				{
					Name: "baz", Broken: true,
					PossibleStates: fsm.MustNewStateMap(
						fsm.State{Name: "lol"},
					),
				},
			},
			TestCase: pl_testing.TestCase{
				MustFail: true,
				MustFailIsErr: fsm.UnexpectedPossibleStatesError{
					StateName:          "baz",
					MustBePresent:      false,
					PossibleStateNames: "lol",
				},
			},
		},
		"case06: success": {
			inputStates: []fsm.State{
				{Name: "baz", Broken: true},
				{Name: "bar", Final: true},
				{
					Name: "foo", Initial: true,
					PossibleStates: fsm.MustNewStateMap(
						fsm.State{Name: "lol"},
					),
				},
			},
			expectedInitialState: fsm.State{Name: "foo"},
			expectedFinalState:   fsm.State{Name: "bar"},
			expectedBrokenState:  fsm.State{Name: "baz"},
		},
	}

	for name, tc := range cases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			smap, err := fsm.NewStateMap(tc.inputStates...)
			require.NoError(t, err)

			tc.CheckError(t,
				smap.Validate(),
			)

			var (
				initialState, initialStateExists = smap.Initial()
				finalState, finalStateExists     = smap.Final()
				brokenState, brokenStateExists   = smap.Broken()
			)

			require.True(t, initialStateExists)
			require.True(t, finalStateExists)
			require.True(t, brokenStateExists)

			require.True(t, tc.expectedInitialState.Equal(initialState))
			require.True(t, tc.expectedFinalState.Equal(finalState))
			require.True(t, tc.expectedBrokenState.Equal(brokenState))
		})
	}
}
