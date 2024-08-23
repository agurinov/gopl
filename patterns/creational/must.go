package creational

import "fmt"

func Must[T any](t T, err error) T {
	if err != nil {
		panic(fmt.Errorf("can't construct %T object\n%w", t, err))
	}

	return t
}
