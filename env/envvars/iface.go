package envvars

import (
	"fmt"
	"net"
	"net/url"
	"time"

	"github.com/google/uuid"
)

type T interface {
	string | bool | int | time.Duration | uuid.UUID | net.IP | url.URL
}

type Variable[V T] interface {
	Present() bool
	Value() (V, error)
	Store(*V) error
	fmt.Stringer
}

func String(key string) Variable[string] {
	return impl[string]{
		key:    key,
		mapper: toStringMapper,
	}
}

func Bool(key string) Variable[bool] {
	return impl[bool]{
		key:    key,
		mapper: toBoolMapper,
	}
}

func Int(key string) Variable[int] {
	return impl[int]{
		key:    key,
		mapper: toIntMapper,
	}
}

func Duration(key string) Variable[time.Duration] {
	return impl[time.Duration]{
		key:    key,
		mapper: toDurationMapper,
	}
}

func UUID(key string) Variable[uuid.UUID] {
	return impl[uuid.UUID]{
		key:    key,
		mapper: toUUIDMapper,
	}
}

func IP(key string) Variable[net.IP] {
	return impl[net.IP]{
		key:    key,
		mapper: toIPMapper,
	}
}

func URL(key string) Variable[url.URL] {
	return impl[url.URL]{
		key:    key,
		mapper: toURLMapper,
	}
}
