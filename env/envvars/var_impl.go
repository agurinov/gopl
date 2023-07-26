package envvars

import (
	"fmt"
	"os"
)

type impl[V variable] struct {
	mapper func(string) (V, error)
	key    string
}

func (i impl[V]) String() string {
	return i.key
}

func (i impl[V]) Present() bool {
	_, present := os.LookupEnv(i.String())

	return present
}

func (i impl[V]) Value() (V, error) {
	var typed V

	value, present := os.LookupEnv(i.String())
	if !present {
		return typed, fmt.Errorf(
			"envvar %q doesn't present",
			i,
		)
	}

	if i.mapper == nil {
		return typed, fmt.Errorf(
			"envvar %q doesn't have mapper for type %T",
			i,
			typed,
		)
	}

	return i.mapper(value)
}

func (i impl[V]) Store(dst *V) error {
	typed, err := i.Value()
	if err != nil {
		return err
	}

	*dst = typed

	return nil
}
