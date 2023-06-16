package fsm

import (
	"context"

	c "github.com/agurinov/gopl/patterns/creational"
)

type StateMachine[C Context] struct {
	states  StateMap
	storage StateStorage[C]
	name    string
	version string
}

func (m StateMachine[C]) currentState(
	ctx context.Context,
	event Event[C],
) (State, error) {
	if m.storage == nil {
		initialState, initialStateExists := m.states.Initial()

		if !initialStateExists {
			return EmptyState, StateNotPresentError{StateName: "initial"}
		}

		return initialState, nil
	}

	currentState, err := m.storage.GetState(ctx, event.Context)
	if err != nil {
		return EmptyState, err
	}

	return currentState, nil
}

func (m StateMachine[C]) isValidTransition(
	currentState State,
	newState State,
) bool {
	if newState.Equal(currentState) {
		return true
	}

	if brokenState, brokenStateExists := m.states.Broken(); brokenStateExists {
		if newState.Equal(brokenState) {
			return true
		}
	}

	if possibleState, possibleExists := currentState.PossibleStates.get(
		newState.Name,
	); possibleExists {
		return newState.Equal(possibleState)
	}

	return false
}

func (m StateMachine[C]) Transition(
	ctx context.Context,
	event Event[C],
) (State, error) {
	if err := event.Validate(); err != nil {
		return EmptyState, err
	}

	currentState, err := m.currentState(ctx, event)
	if err != nil {
		return EmptyState, err
	}

	switch {
	case currentState.Final:
		return currentState, nil
	case currentState.Broken:
		return currentState, nil
	}

	newState, err := event.TransitionFunc(ctx)
	if err != nil {
		return EmptyState, err
	}

	if sameState := newState.Equal(currentState); sameState {
		return newState, nil
	}

	if validTransition := m.isValidTransition(currentState, newState); !validTransition {
		brokenState, brokenStateExists := m.states.Broken()
		if !brokenStateExists {
			return EmptyState, StateNotPresentError{StateName: "broken"}
		}

		newState = brokenState
	}

	if newState.OnTransition != nil {
		if err := newState.OnTransition(ctx); err != nil {
			return EmptyState, err
		}
	}

	if m.storage != nil {
		if err := m.storage.PushState(ctx, event.Context, newState); err != nil {
			return EmptyState, err
		}
	}

	return newState, nil
}

func New[C Context](opts ...StateMachineOption[C]) (StateMachine[C], error) {
	return c.New(opts...)
}
