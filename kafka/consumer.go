package kafka

import (
	"context"
	"fmt"

	"github.com/agurinov/gopl/backoff"
	pl_errors "github.com/agurinov/gopl/errors"
	c "github.com/agurinov/gopl/patterns/creational"
)

type Consumer[E Event] struct {
	library             ConsumerLibrary
	dlq                 ProducerLibrary
	eventBatchHandler   EventBatchHandler[E]
	eventHandler        EventHandler[E]
	backoff             *backoff.Backoff
	configmap           ConfigMap
	eventSerializer     EventSerializer[E]
	cfg                 Config
	maxIterations       uint64
	eventHandleStrategy EventHandleStrategy
}

//nolint:gocognit,gocyclo,cyclop,maintidx,nakedret
func (c Consumer[E]) Consume(ctx context.Context) (err error) {
	if err = c.library.Init(ctx, c.configmap, c.cfg); err != nil {
		err = fmt.Errorf("can't init consumer library: %w", err)

		return
	}

	defer func() {
		if libraryCloseErr := c.library.Close(); libraryCloseErr != nil {
			err = pl_errors.Or(
				err,
				fmt.Errorf("can't close consumer library: %w", libraryCloseErr),
			)
		}
	}()

	if c.dlq != nil {
		if err = c.dlq.Init(ctx, c.configmap, c.cfg); err != nil {
			err = fmt.Errorf("can't init dlq library: %w", err)

			return
		}

		defer func() {
			if dlqCloseErr := c.dlq.Close(); dlqCloseErr != nil {
				err = pl_errors.Or(
					err,
					fmt.Errorf("can't close dlq library: %w", dlqCloseErr),
				)
			}
		}()
	}

	type e struct {
		serialized E
		raw        []byte
	}

	var (
		dlqEvents    = make([][]byte, 0, c.cfg.BatchSize)
		batch        = make([]E, 0, c.cfg.BatchSize)
		batchWrapped = make([]e, 0, c.cfg.BatchSize)
	)

	for i := 0; c.maxIterations == 0 || i < int(c.maxIterations); i++ {
		// Clear slice this way to reuse same underlying array every time.
		dlqEvents = dlqEvents[:0]
		batch = batch[:0]
		batchWrapped = batchWrapped[:0]

		var (
			rawEvents [][]byte
			position  EventPosition
		)

	CONSUME:
		rawEvents, position, err = c.library.ConsumeBatch(ctx, c.cfg.BatchSize)

		if len(rawEvents) == 0 || err != nil {
			if err = healthcheck(ctx, c.backoff, err); err != nil {
				err = fmt.Errorf("can't consume: %w", err)

				return
			}

			goto CONSUME
		}

		c.backoff.Reset()

		if err = position.ValidateWith(c.cfg.EventPosition); err != nil {
			err = fmt.Errorf("can't consume: %w", err)

			return
		}

	SERIALIZE:
		for i := range rawEvents {
			event, serializeErr := c.eventSerializer(rawEvents[i])

			switch {
			case serializeErr == nil:
				batchWrapped = append(batchWrapped, e{
					raw:        rawEvents[i],
					serialized: event,
				})
			case c.dlq != nil:
				err = serializeErr
				dlqEvents = append(dlqEvents, rawEvents[i])
			default:
				err = serializeErr

				break SERIALIZE
			}
		}

		if healtcheckErr := healthcheck(ctx, nil, nil); healtcheckErr != nil {
			err = pl_errors.Or(err, healtcheckErr)

			return
		}

		var handledEvents int

		switch c.eventHandleStrategy {
		case EventHandleBatch:
			// Ignore dlq at this strategy
			// So far we don't know exactly which events from batch determined as broken
			for i := range batchWrapped {
				batch = append(batch, batchWrapped[i].serialized)
			}

			if handleErr := c.eventBatchHandler.Handle(ctx, batch); handleErr != nil {
				err = pl_errors.Or(
					err,
					fmt.Errorf("can't handle event batch: %w", handleErr),
				)

				return
			}

			handledEvents += len(batch)
		case EventHandleOneByOne:
		loop:
			for i := range batchWrapped {
				handleErr := c.eventHandler.Handle(ctx, batchWrapped[i].serialized)

				switch {
				case handleErr == nil:
					handledEvents++
				case c.dlq != nil:
					dlqEvents = append(dlqEvents, batchWrapped[i].raw)
				default:
					err = pl_errors.Or(
						err,
						fmt.Errorf("can't handle event: %w", handleErr),
					)

					break loop
				}
			}
		}

		var (
			totalHandledEvents   = handledEvents + len(dlqEvents)
			isSomethingToProduce = c.dlq != nil && len(dlqEvents) > 0
			isSomethingToCommit  = totalHandledEvents > 0
		)

	PRODUCE_DLQ:
		if isSomethingToProduce {
			if produceErr := c.dlq.ProduceBatch(ctx, dlqEvents...); produceErr != nil {
				if healthcheckErr := healthcheck(ctx, c.backoff, produceErr); healthcheckErr != nil {
					err = pl_errors.Or(
						err,
						fmt.Errorf("can't produce dlq: %w", healthcheckErr),
					)

					return
				}

				goto PRODUCE_DLQ
			}

			c.backoff.Reset()
		}

		position.Offset += int64(totalHandledEvents - 1)

	COMMIT_OFFSET:
		if isSomethingToCommit {
			if commitErr := c.library.Commit(ctx, position); commitErr != nil {
				if healthcheckErr := healthcheck(ctx, c.backoff, commitErr); healthcheckErr != nil {
					err = pl_errors.Or(
						err,
						fmt.Errorf("can't commit event position: %w", healthcheckErr),
					)

					return
				}

				goto COMMIT_OFFSET
			}

			c.backoff.Reset()
		}
	}

	return
}

func (c Consumer[E]) Validate() error {
	if c.library == nil {
		return ErrEmptyConsumerLibrary
	}

	if c.eventSerializer == nil {
		return ErrEmptySerializer
	}

	if err := c.cfg.Validate(); err != nil {
		return err
	}

	if err := c.configmap.ValidateForConsumer(); err != nil {
		return err
	}

	if c.dlq != nil {
		if err := c.configmap.ValidateForProducer(); err != nil {
			return err
		}
	}

	switch c.eventHandleStrategy {
	case EventHandleOneByOne:
		if c.eventHandler == nil {
			return fmt.Errorf(
				"%w: <nil> handler with %s strategy",
				ErrInvalidEventHandler,
				EventHandleOneByOne,
			)
		}
	case EventHandleBatch:
		if c.eventBatchHandler == nil {
			return fmt.Errorf(
				"%w: <nil> batch handler with %s strategy",
				ErrInvalidEventHandler,
				EventHandleBatch,
			)
		}
	default:
		return fmt.Errorf(
			"%w: unsupported handle strategy %s",
			ErrInvalidEventHandler,
			c.eventHandleStrategy,
		)
	}

	return nil
}

func NewConsumer[E Event](opts ...c.Option[Consumer[E]]) (Consumer[E], error) {
	obj, err := c.New(opts...)
	if err != nil {
		return obj, err
	}

	if err := obj.Validate(); err != nil {
		return obj, fmt.Errorf(
			"can't validate consumer: %w",
			err,
		)
	}

	return obj, nil
}
