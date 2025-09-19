package ddd

import (
	"errors"
	"fmt"
)

type (
	ConflictError[T any] struct{}
	conflict             interface{ isConflict() }
)

func (ConflictError[T]) Error() string {
	var zero T

	return fmt.Sprintf("%T: conflict", zero)
}

func (ConflictError[T]) isConflict() {}

func IsConflict(err error) bool {
	var iface conflict

	return errors.As(err, &iface)
}
