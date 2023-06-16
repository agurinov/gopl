package fsm

import (
	"context"
	"fmt"
)

type (
	Context          any
	Event[C Context] struct {
		Context        C
		TransitionFunc func(context.Context) (State, error)
	}
)

func (e Event[C]) Validate() error {
	if e.TransitionFunc == nil {
		return fmt.Errorf("%w: transition function must be present", ErrInvalidEvent)
	}

	return nil
}
