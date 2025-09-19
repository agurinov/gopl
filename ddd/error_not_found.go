package ddd

import (
	"errors"
	"fmt"
)

type (
	NotFoundError[T any] struct{}
	notFound             interface{ isNotFound() }
)

func (NotFoundError[T]) Error() string {
	var zero T

	return fmt.Sprintf("no %T found", zero)
}

func (NotFoundError[T]) isNotFound() {}

func IsNotFound(err error) bool {
	var iface notFound

	return errors.As(err, &iface)
}
