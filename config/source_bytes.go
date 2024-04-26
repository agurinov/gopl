package config

import (
	"context"
)

func FromBytes(content []byte) Source {
	return func(context.Context) ([]byte, error) {
		return content, nil
	}
}
