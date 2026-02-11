package kafka

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

type (
	RecordDiscarder[V any] interface {
		Discard(context.Context, ...V) error
		// Commit(context.Context, ...V) error
	}
	dlqRecordDiscarder[V any] struct {
		dlq Producer[V]
	}
	noopRecordDiscarder[V fmt.Stringer] struct {
		logger *zap.Logger
	}
)

func (d noopRecordDiscarder[V]) Discard(
	_ context.Context,
	records ...V,
) error {
	d.logger.Warn(
		"discard records",
		zap.String("action", "noop"),
		zap.Stringers("records", records),
	)

	return nil
}

func (d dlqRecordDiscarder[V]) Discard(
	ctx context.Context,
	records ...V,
) error {
	return d.dlq.Produce(ctx, records...)
}
