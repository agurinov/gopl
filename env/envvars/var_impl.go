package pl_envvars

import (
	"io"
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

	// TODO(a.gurinov): Fix errors
	value, present := os.LookupEnv(i.String())
	if !present {
		return typed, io.EOF
	}

	if i.mapper == nil {
		return typed, io.EOF
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
