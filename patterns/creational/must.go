package creational

import "fmt"

func Must[T any](t T, err error) T {
	if err != nil {
		panic(fmt.Errorf("can't construct %T object\n%w", t, err))
	}

	return t
}

func MustInit(err error) {
	if err != nil {
		panic(fmt.Errorf("can't init object\n%w", err))
	}
}
