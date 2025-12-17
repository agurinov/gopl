package kafka

import (
	"context"

	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
)

type (
	RecordDiscarder[V any] interface {
		Discard(context.Context, ...V) error
		// Commit(context.Context, ...V) error
	}
	dlqRecordDiscarder  struct{}
	noopRecordDiscarder struct {
		logger *zap.Logger
	}
)

func (d noopRecordDiscarder) Discard(
	_ context.Context,
	records ...*kgo.Record,
) error {
	return nil
}

func (d dlqRecordDiscarder) Discard(
	_ context.Context,
	records ...*kgo.Record,
) error {
	return nil
}
