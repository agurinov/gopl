package kafka

import (
	"context"
	"fmt"

	"go.uber.org/zap"
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
	d.logger.Warn(
		"discard records to log",
		// RecordLogFields()...,
		// zap.Objects("records", records),
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
