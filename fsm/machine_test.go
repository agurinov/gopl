//go:build test_unit

package fsm_test

import (
	"context"
	"io"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/agurinov/gopl/fsm"
	pl_testing "github.com/agurinov/gopl/testing"
	pl_gomock "github.com/agurinov/gopl/testing/gomock"
)

func TestStateMachine_Transition(t *testing.T) {
	pl_testing.Init(t)

	registrationContext := RegistrationContext{
		UUID: uuid.MustParse("10000000-0000-0000-0000-111111111111"),
	}

	cases := map[string]struct {
		mocks         func(mocks)
		inputOptions  []RegistrationStateMachineOption
		inputEvent    RegistrationEvent
		expectedState fsm.State
		pl_testing.TestCase
	}{
		"case00: invalid event: without transition func": {
			inputEvent: RegistrationEvent{
				Context:        registrationContext,
				TransitionFunc: nil,
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: fsm.ErrInvalidEvent,
			},
		},
		"case01: stateless machine err on transition": {
			inputOptions: []RegistrationStateMachineOption{
				fsm.WithStateStorage[RegistrationContext](nil),
			},
			inputEvent: RegistrationEvent{
				Context: registrationContext,
				TransitionFunc: func(_ context.Context) (fsm.State, error) {
					return approvedState, io.EOF
				},
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: io.EOF,
			},
		},
		"case02: stateless machine unexpected transition": {
			inputOptions: []RegistrationStateMachineOption{
				fsm.WithStateStorage[RegistrationContext](nil),
			},
			inputEvent: RegistrationEvent{
				Context: registrationContext,
				TransitionFunc: func(_ context.Context) (fsm.State, error) {
					return approvedState, nil
				},
			},
			expectedState: brokenState,
		},
		"case03: stateless machine explicit broken transition": {
			inputOptions: []RegistrationStateMachineOption{
				fsm.WithStateStorage[RegistrationContext](nil),
			},
			inputEvent: RegistrationEvent{
				Context: registrationContext,
				TransitionFunc: func(_ context.Context) (fsm.State, error) {
					return brokenState, nil
				},
			},
			expectedState: brokenState,
		},
		"case04: stateless machine expected transition": {
			inputOptions: []RegistrationStateMachineOption{
				fsm.WithStateStorage[RegistrationContext](nil),
			},
			inputEvent: RegistrationEvent{
				Context: registrationContext,
				TransitionFunc: func(_ context.Context) (fsm.State, error) {
					return uploadPassportState, nil
				},
			},
			expectedState: uploadPassportState,
		},

		"case05: stateful machine err on current state": {
			mocks: func(m mocks) {
				gomock.InOrder(
					m.storage.EXPECT().
						GetState(pl_gomock.IsContext(), registrationContext).
						Times(1).
						Return(fsm.EmptyState, io.EOF),
				)
			},
			inputEvent: RegistrationEvent{
				Context: registrationContext,
				TransitionFunc: func(_ context.Context) (fsm.State, error) {
					return uploadPassportState, nil
				},
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: io.EOF,
			},
		},
		"case06: stateful machine in final state": {
			mocks: func(m mocks) {
				gomock.InOrder(
					m.storage.EXPECT().
						GetState(pl_gomock.IsContext(), registrationContext).
						Times(1).
						Return(approvedState, nil),
				)
			},
			inputEvent: RegistrationEvent{
				Context: registrationContext,
				TransitionFunc: func(_ context.Context) (fsm.State, error) {
					return uploadPassportState, nil
				},
			},
			expectedState: approvedState,
		},
		"case07: stateful machine in broken state": {
			mocks: func(m mocks) {
				gomock.InOrder(
					m.storage.EXPECT().
						GetState(pl_gomock.IsContext(), registrationContext).
						Times(1).
						Return(brokenState, nil),
				)
			},
			inputEvent: RegistrationEvent{
				Context: registrationContext,
				TransitionFunc: func(_ context.Context) (fsm.State, error) {
					return uploadPassportState, nil
				},
			},
			expectedState: brokenState,
		},
		"case08: stateful machine new state transition err": {
			mocks: func(m mocks) {
				gomock.InOrder(
					m.storage.EXPECT().
						GetState(pl_gomock.IsContext(), registrationContext).
						Times(1).
						Return(chooseCountryState, nil),
				)
			},
			inputEvent: RegistrationEvent{
				Context: registrationContext,
				TransitionFunc: func(_ context.Context) (fsm.State, error) {
					return chooseCountryState, io.EOF
				},
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: io.EOF,
			},
		},
		"case09: stateful machine same state idempotence check": {
			mocks: func(m mocks) {
				gomock.InOrder(
					m.storage.EXPECT().
						GetState(pl_gomock.IsContext(), registrationContext).
						Times(1).
						Return(deniedState, nil),
				)
			},
			inputEvent: RegistrationEvent{
				Context: registrationContext,
				TransitionFunc: func(_ context.Context) (fsm.State, error) {
					return deniedState, nil
				},
			},
			expectedState: deniedState,
		},
		"case10: stateful machine unexpected transition": {
			mocks: func(m mocks) {
				gomock.InOrder(
					m.storage.EXPECT().
						GetState(pl_gomock.IsContext(), registrationContext).
						Times(1).
						Return(uploadDriverLicenseState, nil),
					m.storage.EXPECT().
						PushState(
							pl_gomock.IsContext(),
							registrationContext,
							pl_gomock.Eq(brokenState),
						).
						Times(1).
						Return(nil),
				)
			},
			inputEvent: RegistrationEvent{
				Context: registrationContext,
				TransitionFunc: func(_ context.Context) (fsm.State, error) {
					return chooseCountryState, nil
				},
			},
			expectedState: brokenState,
		},
		"case11: stateful machine expected transition": {
			mocks: func(m mocks) {
				gomock.InOrder(
					m.storage.EXPECT().
						GetState(pl_gomock.IsContext(), registrationContext).
						Times(1).
						Return(uploadDriverLicenseState, nil),
					m.storage.EXPECT().
						PushState(
							pl_gomock.IsContext(),
							registrationContext,
							reviewState,
						).
						Times(1).
						Return(nil),
				)
			},
			inputEvent: RegistrationEvent{
				Context: registrationContext,
				TransitionFunc: func(_ context.Context) (fsm.State, error) {
					return reviewState, nil
				},
			},
			expectedState: reviewState,
		},
		"case12: stateful machine err on push state": {
			mocks: func(m mocks) {
				gomock.InOrder(
					m.storage.EXPECT().
						GetState(pl_gomock.IsContext(), registrationContext).
						Times(1).
						Return(uploadPassportState, nil),
					m.storage.EXPECT().
						PushState(
							pl_gomock.IsContext(),
							registrationContext,
							pl_gomock.Eq(uploadSelfieState),
						).
						Times(1).
						Return(io.EOF),
				)
			},
			inputEvent: RegistrationEvent{
				Context: registrationContext,
				TransitionFunc: func(_ context.Context) (fsm.State, error) {
					return uploadSelfieState, nil
				},
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: io.EOF,
			},
		},
	}

	for name, tc := range cases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			ctx := context.TODO()

			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			m := NewMocks(ctrl)
			if tc.mocks != nil {
				tc.mocks(m)
			}

			opts := []RegistrationStateMachineOption{
				fsm.WithName[RegistrationContext]("registration_machine"),
				fsm.WithVersion[RegistrationContext]("v1"),
				fsm.WithStateStorage[RegistrationContext](m.storage),
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
			opts = append(opts, tc.inputOptions...)

			sm, err := fsm.New(opts...)
			require.NoError(t, err)

			state, err := sm.Transition(ctx, tc.inputEvent)

			tc.CheckError(t, err)
			require.True(t, tc.expectedState.Equal(state))
		})
	}
}
