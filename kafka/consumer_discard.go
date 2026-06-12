package kafka

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/agurinov/gopl/x"
)

type (
	RecordDiscarder interface {
		Discard(context.Context, ...Record) error
	}
	dlqRecordDiscarder struct {
		dlq Producer
	}
	logRecordDiscarder struct {
		logger *zap.Logger
	}
)

func (d logRecordDiscarder) Discard(
	_ context.Context,
	records ...Record,
) error {
	loggableRecords := x.SliceConvert(
		records,
		func(in Record) loggableRecord {
			return loggableRecord{
				Record: in,
			}
		},
	)

	d.logger.Warn(
		"discard records to log",
		zap.Objects("records", loggableRecords),
	)

	return nil
}

func (d dlqRecordDiscarder) Discard(
	ctx context.Context,
	records ...Record,
) error {
	if err := d.dlq.Produce(ctx, records...); err != nil {
		return fmt.Errorf("can't dlq: %w", err)
	}

	return nil
}
