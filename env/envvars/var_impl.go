package pl_envvars

import (
	"io"
	"os"
)

type impl[V variable] struct {
	mapper func(string) (V, error)
	key    string
}

func (v impl[V]) String() string {
	return v.key
}

func (v impl[V]) Value() (V, error) {
	var typed V

	value, exists := os.LookupEnv(v.String())
	if !exists {
		return typed, io.EOF
	}

	if v.mapper == nil {
		return typed, io.EOF
	}

	return v.mapper(value)
}

func (v impl[V]) Store(dst *V) error {
	typed, err := v.Value()
	if err != nil {
		return err
	}

	*dst = typed

	return nil
}
