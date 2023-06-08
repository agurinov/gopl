package pl_testing

import "testing"

type (
	Stand interface {
		Up(*testing.T) bool
		Name() StandName
	}
	StandName  string
	StandState map[StandName]struct{}
)
