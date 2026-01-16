package fsm

import "context"

type StateStorage[C Context] interface {
	GetState(context.Context, C) (State, error)
	PushState(context.Context, C, State) error
}
