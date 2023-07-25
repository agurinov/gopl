package fsm

import (
	c "github.com/agurinov/gopl/patterns/creational"
)

type StateMachineOption[C Context] c.Option[StateMachine[C]]

func WithName[C Context](name string) StateMachineOption[C] {
	return func(sm *StateMachine[C]) error {
		sm.name = name

		return nil
	}
}

func WithVersion[C Context](version string) StateMachineOption[C] {
	return func(sm *StateMachine[C]) error {
		sm.version = version

		return nil
	}
}

func WithStateStorage[C Context](storage StateStorage[C]) StateMachineOption[C] {
	return func(sm *StateMachine[C]) error {
		sm.storage = storage

		return nil
	}
}

func WithStateMap[C Context](states ...State) StateMachineOption[C] {
	return func(sm *StateMachine[C]) error {
		statemap, err := NewStateMap(states...)
		if err != nil {
			return err
		}

		if validateErr := statemap.Validate(); validateErr != nil {
			return validateErr
		}

		sm.states = statemap

		return nil
	}
}
