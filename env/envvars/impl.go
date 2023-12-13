package envvars

import (
	"fmt"
	"os"
)

type impl[V T] struct {
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
		return typed, fmt.Errorf("%q: %w", i, ErrNoVar)
	}

	if i.mapper == nil {
		return typed, fmt.Errorf("%q: %w for type %T", i, ErrNoMapper, typed)
	}

	typed, err := i.mapper(value)
	if err != nil {
		return typed, fmt.Errorf("%q: %w", i, err)
	}

	return typed, nil
}

func (i impl[V]) Store(dst *V) error {
	typed, err := i.Value()
	if err != nil {
		return err
	}

	*dst = typed

	return nil
}
