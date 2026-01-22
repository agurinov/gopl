package kafka

import "context"

type (
	Producer[V any] interface {
		Produce(context.Context, ...V) error
	}
)
