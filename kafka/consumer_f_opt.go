package kafka

import (
	"github.com/agurinov/gopl/backoff"
	c "github.com/agurinov/gopl/patterns/creational"
)

type ConsumerOption[E Event] c.Option[Consumer[E]]

func WithLibrary[E Event](library ConsumerLibrary) ConsumerOption[E] {
	return func(consumer *Consumer[E]) error {
		consumer.library = library

		return nil
	}
}

func WithDLQ[E Event](dlq ProducerLibrary) ConsumerOption[E] {
	return func(consumer *Consumer[E]) error {
		consumer.dlq = dlq

		return nil
	}
}

func WithConfigMap[E Event](configmap ConfigMap) ConsumerOption[E] {
	return func(consumer *Consumer[E]) error {
		consumer.configmap = configmap

		return nil
	}
}

func WithConsumerConfig[E Event](cfg Config) ConsumerOption[E] {
	return func(consumer *Consumer[E]) error {
		consumer.cfg = cfg

		return nil
	}
}

func WithEventSerializer[E Event](s EventSerializer[E]) ConsumerOption[E] {
	return func(consumer *Consumer[E]) error {
		consumer.eventSerializer = s

		return nil
	}
}

func WithEventHandleStrategy[E Event](s EventHandleStrategy) ConsumerOption[E] {
	return func(consumer *Consumer[E]) error {
		consumer.eventHandleStrategy = s

		return nil
	}
}

func WithEventHandler[E Event](h EventHandler[E]) ConsumerOption[E] {
	return func(consumer *Consumer[E]) error {
		consumer.eventHandleStrategy = EventHandleOneByOne
		consumer.eventHandler = h

		return nil
	}
}

func WithEventBatchHandler[E Event](h EventBatchHandler[E]) ConsumerOption[E] {
	return func(consumer *Consumer[E]) error {
		consumer.eventHandleStrategy = EventHandleBatch
		consumer.eventBatchHandler = h

		return nil
	}
}

func WithMaxIterations[E Event](i uint64) ConsumerOption[E] {
	return func(consumer *Consumer[E]) error {
		consumer.maxIterations = i

		return nil
	}
}

func WithBackoffOptions[E Event](opts ...backoff.BackoffOption) ConsumerOption[E] {
	return func(consumer *Consumer[E]) error {
		opts = append(opts,
			backoff.WithName("kafka-consumer"),
		)

		b, err := backoff.New(opts...)
		if err != nil {
			return err
		}

		consumer.backoff = b

		return nil
	}
}
