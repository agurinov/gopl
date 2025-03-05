package creational

import "fmt"

func MustInit(err error) {
	if err != nil {
		panic(fmt.Errorf("can't init object\n%w", err))
	}
}

func Must[T any](t T, err error) T {
	if err != nil {
		panic(fmt.Errorf("can't construct %T object\n%w", t, err))
	}

	return t
}

func MustDuo[T1 any, T2 any](t1 T1, t2 T2, err error) (T1, T2) {
	if err != nil {
		panic(fmt.Errorf("can't construct %T and %T objects\n%w", t1, t2, err))
	}

	return t1, t2
}
