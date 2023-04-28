package kafka

import (
	"github.com/agurinov/gopl/backoff"
	c "github.com/agurinov/gopl/patterns/creational"
)

func WithLibrary[E Event](library ConsumerLibrary) c.Option[Consumer[E]] {
	return func(consumer *Consumer[E]) error {
		consumer.library = library

		return nil
	}
}

func WithDLQ[E Event](dlq ProducerLibrary) c.Option[Consumer[E]] {
	return func(consumer *Consumer[E]) error {
		consumer.dlq = dlq

		return nil
	}
}

func WithConfigMap[E Event](configmap ConfigMap) c.Option[Consumer[E]] {
	return func(consumer *Consumer[E]) error {
		consumer.configmap = configmap

		return nil
	}
}

func WithConsumerConfig[E Event](cfg Config) c.Option[Consumer[E]] {
	return func(consumer *Consumer[E]) error {
		consumer.cfg = cfg

		return nil
	}
}

func WithEventSerializer[E Event](s EventSerializer[E]) c.Option[Consumer[E]] {
	return func(consumer *Consumer[E]) error {
		consumer.eventSerializer = s

		return nil
	}
}

func WithEventHandleStrategy[E Event](s EventHandleStrategy) c.Option[Consumer[E]] {
	return func(consumer *Consumer[E]) error {
		consumer.eventHandleStrategy = s

		return nil
	}
}

func WithEventHandler[E Event](h EventHandler[E]) c.Option[Consumer[E]] {
	return func(consumer *Consumer[E]) error {
		consumer.eventHandleStrategy = EventHandleOneByOne
		consumer.eventHandler = h

		return nil
	}
}

func WithEventBatchHandler[E Event](h EventBatchHandler[E]) c.Option[Consumer[E]] {
	return func(consumer *Consumer[E]) error {
		consumer.eventHandleStrategy = EventHandleBatch
		consumer.eventBatchHandler = h

		return nil
	}
}

func WithMaxIterations[E Event](i uint64) c.Option[Consumer[E]] {
	return func(consumer *Consumer[E]) error {
		consumer.maxIterations = i

		return nil
	}
}

func WithBackoffOptions[E Event](opts ...backoff.BackoffOption) c.Option[Consumer[E]] {
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
