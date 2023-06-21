package fsm

import (
	"context"
)

type State struct {
	OnTransition   func(context.Context) error
	PossibleStates StateMap
	Name           string
	Initial        bool
	Final          bool
	Broken         bool
}

var EmptyState = State{}

func (s State) Equal(other State) bool {
	return s.Name == other.Name
}
