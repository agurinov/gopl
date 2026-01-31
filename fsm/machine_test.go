package fsm_test

import (
	"context"
	"io"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/agurinov/gopl/fsm"
	"github.com/agurinov/gopl/fsm/mockery"
	pl_testing "github.com/agurinov/gopl/testing"
)

func TestStateMachine_Transition(t *testing.T) {
	pl_testing.Init(t)

	type (
		di struct {
			storage *mockery.MockStateStorage[RegistrationContext]
		}
		args struct {
			event   RegistrationEvent
			options []RegistrationStateMachineOption
		}
		results struct {
			state fsm.State
		}
	)

	registrationContext := RegistrationContext{
		UUID: uuid.MustParse("10000000-0000-0000-0000-111111111111"),
	}

	cases := map[string]struct {
		args args
		di   func(*testing.T, di)
		pl_testing.TestCase
		results results
	}{
		"case00: invalid event: without transition func": {
			args: args{
				event: RegistrationEvent{
					Context:        registrationContext,
					TransitionFunc: nil,
				},
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: fsm.ErrInvalidEvent,
			},
		},
		"case01: stateless machine err on transition": {
			args: args{
				options: []RegistrationStateMachineOption{
					fsm.WithStateStorage[RegistrationContext](nil),
				},
				event: RegistrationEvent{
					Context: registrationContext,
					TransitionFunc: func(_ context.Context) (fsm.State, error) {
						return approvedState, io.EOF
					},
				},
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: io.EOF,
			},
		},
		"case02: stateless machine unexpected transition": {
			args: args{
				options: []RegistrationStateMachineOption{
					fsm.WithStateStorage[RegistrationContext](nil),
				},
				event: RegistrationEvent{
					Context: registrationContext,
					TransitionFunc: func(_ context.Context) (fsm.State, error) {
						return approvedState, nil
					},
				},
			},
			results: results{
				state: brokenState,
			},
		},
		"case03: stateless machine explicit broken transition": {
			args: args{
				options: []RegistrationStateMachineOption{
					fsm.WithStateStorage[RegistrationContext](nil),
				},
				event: RegistrationEvent{
					Context: registrationContext,
					TransitionFunc: func(_ context.Context) (fsm.State, error) {
						return brokenState, nil
					},
				},
			},
			results: results{
				state: brokenState,
			},
		},
		"case04: stateless machine expected transition": {
			args: args{
				options: []RegistrationStateMachineOption{
					fsm.WithStateStorage[RegistrationContext](nil),
				},
				event: RegistrationEvent{
					Context: registrationContext,
					TransitionFunc: func(_ context.Context) (fsm.State, error) {
						return uploadPassportState, nil
					},
				},
			},
			results: results{
				state: uploadPassportState,
			},
		},
		"case05: stateful machine err on current state": {
			di: func(t *testing.T, d di) {
				t.Helper()

				d.storage.
					On(
						"GetState",
						mock.AnythingOfType("*context.cancelCtx"),
						registrationContext,
					).
					Return(fsm.EmptyState, io.EOF).
					Times(1)
			},
			args: args{
				event: RegistrationEvent{
					Context: registrationContext,
					TransitionFunc: func(_ context.Context) (fsm.State, error) {
						return uploadPassportState, nil
					},
				},
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: io.EOF,
			},
		},
		"case06: stateful machine in final state": {
			di: func(t *testing.T, d di) {
				t.Helper()

				d.storage.
					On(
						"GetState",
						mock.AnythingOfType("*context.cancelCtx"),
						registrationContext,
					).
					Return(approvedState, nil).
					Times(1)
			},
			args: args{
				event: RegistrationEvent{
					Context: registrationContext,
					TransitionFunc: func(_ context.Context) (fsm.State, error) {
						return uploadPassportState, nil
					},
				},
			},
			results: results{
				state: approvedState,
			},
		},
		"case07: stateful machine in broken state": {
			di: func(t *testing.T, d di) {
				t.Helper()

				d.storage.
					On(
						"GetState",
						mock.AnythingOfType("*context.cancelCtx"),
						registrationContext,
					).
					Return(brokenState, nil).
					Times(1)
			},
			args: args{
				event: RegistrationEvent{
					Context: registrationContext,
					TransitionFunc: func(_ context.Context) (fsm.State, error) {
						return uploadPassportState, nil
					},
				},
			},
			results: results{
				state: brokenState,
			},
		},
		"case08: stateful machine new state transition err": {
			di: func(t *testing.T, d di) {
				t.Helper()

				d.storage.
					On(
						"GetState",
						mock.AnythingOfType("*context.cancelCtx"),
						registrationContext,
					).
					Return(chooseCountryState, nil).
					Times(1)
			},
			args: args{
				event: RegistrationEvent{
					Context: registrationContext,
					TransitionFunc: func(_ context.Context) (fsm.State, error) {
						return chooseCountryState, io.EOF
					},
				},
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: io.EOF,
			},
		},
		"case09: stateful machine same state idempotence check": {
			di: func(t *testing.T, d di) {
				t.Helper()

				d.storage.
					On(
						"GetState",
						mock.AnythingOfType("*context.cancelCtx"),
						registrationContext,
					).
					Return(deniedState, nil).
					Times(1)
			},
			args: args{
				event: RegistrationEvent{
					Context: registrationContext,
					TransitionFunc: func(_ context.Context) (fsm.State, error) {
						return deniedState, nil
					},
				},
			},
			results: results{
				state: deniedState,
			},
		},
		"case10: stateful machine unexpected transition": {
			di: func(t *testing.T, d di) {
				t.Helper()

				d.storage.
					On(
						"GetState",
						mock.AnythingOfType("*context.cancelCtx"),
						registrationContext,
					).
					Return(uploadDriverLicenseState, nil).
					Times(1)

				d.storage.
					On(
						"PushState",
						mock.AnythingOfType("*context.cancelCtx"),
						registrationContext,
						mock.MatchedBy(func(s fsm.State) bool {
							return s.Equal(brokenState)
						}),
					).
					Return(nil).
					Times(1)
			},
			args: args{
				event: RegistrationEvent{
					Context: registrationContext,
					TransitionFunc: func(_ context.Context) (fsm.State, error) {
						return chooseCountryState, nil
					},
				},
			},
			results: results{
				state: brokenState,
			},
		},
		"case11: stateful machine expected transition": {
			di: func(t *testing.T, d di) {
				t.Helper()

				d.storage.
					On(
						"GetState",
						mock.AnythingOfType("*context.cancelCtx"),
						registrationContext,
					).
					Return(uploadDriverLicenseState, nil).
					Times(1)

				d.storage.
					On(
						"PushState",
						mock.AnythingOfType("*context.cancelCtx"),
						registrationContext,
						reviewState,
					).
					Return(nil).
					Times(1)
			},
			args: args{
				event: RegistrationEvent{
					Context: registrationContext,
					TransitionFunc: func(_ context.Context) (fsm.State, error) {
						return reviewState, nil
					},
				},
			},
			results: results{
				state: reviewState,
			},
		},
		"case12: stateful machine err on push state": {
			di: func(t *testing.T, d di) {
				t.Helper()

				d.storage.
					On(
						"GetState",
						mock.AnythingOfType("*context.cancelCtx"),
						registrationContext,
					).
					Return(uploadPassportState, nil).
					Times(1)

				d.storage.
					On(
						"PushState",
						mock.AnythingOfType("*context.cancelCtx"),
						registrationContext,
						uploadSelfieState,
					).
					Return(io.EOF).
					Times(1)
			},
			args: args{
				event: RegistrationEvent{
					Context: registrationContext,
					TransitionFunc: func(_ context.Context) (fsm.State, error) {
						return uploadSelfieState, nil
					},
				},
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: io.EOF,
			},
		},
	}

	for name := range cases {
		tc := cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			diContainer := di{
				storage: mockery.NewMockStateStorage[RegistrationContext](t),
			}

			if tc.di != nil {
				tc.di(t, diContainer)
			}

			ctx := t.Context()

			opts := []RegistrationStateMachineOption{
				fsm.WithName[RegistrationContext]("registration_machine"),
				fsm.WithVersion[RegistrationContext]("v1"),
				fsm.WithStateStorage(diContainer.storage),
				fsm.WithStateMap[RegistrationContext](
					brokenState,
					chooseCountryState,
					uploadPassportState,
					uploadSelfieState,
					uploadDriverLicenseState,
					reviewState,
					deniedState,
					approvedState,
				),
			}
			opts = append(opts, tc.args.options...)

			sm, err := fsm.New(opts...)
			require.NoError(t, err)
			require.NotNil(t, sm)

			state, err := sm.Transition(ctx, tc.args.event)
			tc.CheckError(t, err)
			require.NotNil(t, state)
			require.True(t, tc.results.state.Equal(state))
		})
	}
}
