package kafka

import (
	"context"
)

type (
	Producer interface {
		Produce(context.Context, ...Record) error
		Close(context.Context) error
		Ping(context.Context) error
	}
)
