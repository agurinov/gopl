package kafka

import (
	"context"
	"errors"
	"fmt"

	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"

	"github.com/agurinov/gopl/graceful"
	c "github.com/agurinov/gopl/patterns/creational"
	"github.com/agurinov/gopl/run"
	"github.com/agurinov/gopl/x"
)

type (
	Producer interface {
		Produce(context.Context, ...Record) error
		Close(context.Context) error
	}
	producer struct {
		metrics      producerMetrics
		recordMapper RecordMapper
		logger       *zap.Logger
		client       *kgo.Client
		topic        string
	}
	ProducerOption c.Option[producer]
)

func (p producer) Produce(
	ctx context.Context,
	records ...Record,
) error {
	kgoRecords := x.SliceConvert(
		records,
		p.recordMapper.ToVendor,
	)

	// TODO: force replace topic

	result := p.client.ProduceSync(ctx, kgoRecords...)

	errs := x.SliceConvert(
		result,
		func(in kgo.ProduceResult) error {
			return in.Err
		},
	)

	if err := errors.Join(errs...); err != nil {
		return fmt.Errorf("can't produce records: %w", errors.Join(errs...))
	}

	return nil
}

func (p producer) Close(ctx context.Context) error {
	closeFn := run.ErrorFn(func() error {
		defer p.client.Close()

		if err := p.client.Flush(ctx); err != nil {
			return fmt.Errorf("can't flush records: %w", err)
		}

		return nil
	})

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
