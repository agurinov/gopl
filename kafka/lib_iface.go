package kafka

import (
	"context"
	"io"
)

type ConsumerLibrary interface {
	Init(context.Context, ConfigMap, Config) error
	ConsumeBatch(context.Context, uint) ([][]byte, EventPosition, error)
	Commit(context.Context, EventPosition) error
	io.Closer
}

type ProducerLibrary interface {
	Init(context.Context, ConfigMap, Config) error
	ProduceBatch(context.Context, ...[]byte) error
	io.Closer
}

type ComboLibrary interface {
	ConsumerLibrary
	ProducerLibrary
}
