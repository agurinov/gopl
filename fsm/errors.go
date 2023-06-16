package fsm

import (
	"errors"
	"fmt"
)

var ErrInvalidEvent = errors.New("event: invalid")

type StateAlreadyPresentError struct {
	StateName string
}

func (e StateAlreadyPresentError) Error() string {
	return fmt.Sprintf(
		"statemap: state %q already present",
		e.StateName,
	)
}

type StateNotPresentError struct {
	StateName string
}

func (e StateNotPresentError) Error() string {
	return fmt.Sprintf(
		"statemap: state %q must be present",
		e.StateName,
	)
}

type UnexpectedPossibleStatesError struct {
	PossibleStateNames string
	StateName          string
	MustBePresent      bool
}

func (e UnexpectedPossibleStatesError) Error() string {
	presence := "should not be present"

	if e.MustBePresent {
		presence = "must be present"
	}

	return fmt.Sprintf(
		"statemap: state %q has possible states [%s] which %s",
		e.StateName,
		e.PossibleStateNames,
		presence,
	)
}
