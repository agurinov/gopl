package kafka

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"

	"github.com/agurinov/gopl/graceful"
	c "github.com/agurinov/gopl/patterns/creational"
	"github.com/agurinov/gopl/x"
)

type (
	Producer interface {
		Produce(context.Context, ...Record) error
		Close(context.Context) error
		Ping(context.Context) error
	}
	producer struct {
		metrics      producerMetrics
		recordMapper RecordMapper
		logger       *zap.Logger
		client       *kgo.Client
		topic        string
		cooldown     time.Duration
	}
	ProducerOption c.Option[producer]
)

func (p producer) Produce(
	ctx context.Context,
	records ...Record,
) error {
	kgoRecords, err := x.SliceConvertError(
		records,
		func(in Record) (*kgo.Record, error) {
			in.Topic = cmp.Or(p.topic, in.Topic)

			if p.cooldown > 0 {
				processAfter := time.Now().Add(p.cooldown)

				processAfterBytes, mErr := processAfter.MarshalBinary()
				if mErr != nil {
					return nil, mErr
				}

				in.Headers = append(in.Headers,
					kgo.RecordHeader{
						Key:   cooldownHeaderKey,
						Value: processAfterBytes,
					},
				)
			}

			return p.recordMapper.ToVendor(in), nil
		},
	)
	if err != nil {
		return fmt.Errorf("can't construct records: %w", err)
	}

	result := p.client.ProduceSync(ctx, kgoRecords...)

	errs := x.SliceConvert(
		result,
		func(in kgo.ProduceResult) error {
			return in.Err
		},
	)

	if jErr := errors.Join(errs...); jErr != nil {
		return fmt.Errorf("can't produce records: %w", jErr)
	}

	return nil
}

func (p producer) Close(ctx context.Context) error {
	closeFn := func(ctx context.Context) error {
		defer p.client.Close()

		if err := p.client.Flush(ctx); err != nil {
			return fmt.Errorf("can't flush records: %w", err)
		}

		return nil
	}

	return graceful.Close(closeFn)(ctx)
}

func (p producer) Ping(ctx context.Context) error {
	return p.client.Ping(ctx)
}

func NewProducer(opts ...ProducerOption) (Producer, error) {
	obj := producer{
		recordMapper: kgoRecordMapper{},
	}

	return c.ConstructWithValidate(obj, opts...)
}
