package fsm

import (
	"errors"
	"fmt"
)

type (
	StateAlreadyPresentError struct {
		StateName string
	}
	StateNotPresentError struct {
		StateName string
	}
	UnexpectedPossibleStatesError struct {
		PossibleStateNames string
		StateName          string
		MustBePresent      bool
	}
)

var ErrInvalidEvent = errors.New("event: invalid")

func (e StateAlreadyPresentError) Error() string {
	return fmt.Sprintf(
		"statemap: state %q already present",
		e.StateName,
	)
}

func (e StateNotPresentError) Error() string {
	return fmt.Sprintf(
		"statemap: state %q must be present",
		e.StateName,
	)
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
