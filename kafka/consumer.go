package kafka

import (
	"context"
)

type (
	Consumer interface {
		Start() error
		Close(context.Context) error
		Ping(context.Context) error
	}
)
