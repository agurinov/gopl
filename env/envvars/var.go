package pl_envvars

// TODO(a.gurinov): Errors for this package

import (
	"io"
	"net"
	"os"
	"time"

	"github.com/google/uuid"
)

type variable interface {
	string | bool | int | time.Duration | uuid.UUID | net.IP
}

type Variable[V variable] struct {
	mapper func(string) (V, error)
	key    string
}

func (v Variable[V]) String() string {
	return v.key
}

func (v Variable[V]) Typed() (V, error) {
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

func (v Variable[V]) Store(dst *V) error {
	typed, err := v.Typed()
	if err != nil {
		return err
	}

	*dst = typed

	return nil
}

func String(key string) Variable[string] {
	return Variable[string]{
		key:    key,
		mapper: toStringMapper,
	}
}

func Bool(key string) Variable[bool] {
	return Variable[bool]{
		key:    key,
		mapper: toBoolMapper,
	}
}

func Int(key string) Variable[int] {
	return Variable[int]{
		key:    key,
		mapper: toIntMapper,
	}
}

func Duration(key string) Variable[time.Duration] {
	return Variable[time.Duration]{
		key:    key,
		mapper: toDurationMapper,
	}
}

func UUID(key string) Variable[uuid.UUID] {
	return Variable[uuid.UUID]{
		key:    key,
		mapper: toUUIDMapper,
	}
}

func IP(key string) Variable[net.IP] {
	return Variable[net.IP]{
		key:    key,
		mapper: toIPMapper,
	}
}
