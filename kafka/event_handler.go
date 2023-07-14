package kafka

import "context"

type EventHandleStrategy uint8

const (
	EventHandleOneByOne EventHandleStrategy = iota
	EventHandleBatch
)

type (
	EventHandler[E Event] interface {
		Handle(context.Context, E) error
	}
	EventHandlerFunc[E Event] func(context.Context, E) error
)

func (f EventHandlerFunc[E]) Handle(ctx context.Context, event E) error {
	return f(ctx, event)
}

type (
	EventBatchHandler[E Event] interface {
		Handle(context.Context, []E) error
	}
	EventBatchHandlerFunc[E Event] func(context.Context, []E) error
)

func (f EventBatchHandlerFunc[E]) Handle(ctx context.Context, events []E) error {
	return f(ctx, events)
}
