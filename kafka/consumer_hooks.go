package kafka

import (
	"context"

	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"

	"github.com/agurinov/gopl/internal/x"
)

type consumerHooks struct {
	logger            *zap.Logger
	fetchCounters     *x.Counter[string]
	partitionCounters *x.Counter[int32]
}

func (h consumerHooks) OnFetchRecordBuffered(record Record) {
	key := string(record.Key)

	if h.fetchCounters.Inc(key) > 1 {
		h.logger.Warn(
			"record buffered multiple times",
			RecordLogFields(record)...,
		)

		panic("record already buffered")
	}
}

func (cs consumer) onAssigned(
	ctx context.Context,
	_ *kgo.Client,
	topicPartitions map[string][]int32,
) {
	for topic, assignedPartitions := range topicPartitions {
		if topic != cs.topic {
			cs.logger.Warn(
				"can't assign partitions: not my topic",
				zap.String("got.topic", topic),
				zap.String("expected.topic", cs.topic),
			)

			continue
		}

		// It is VERY important to register the whole batch atomically.
		cs.partitionDispatcher.Run(
			ctx,
			cs.partitionLoop,
			assignedPartitions...,
		)

		cs.logger.Info(
			"kafka partitions assigned",
			zap.String("topic", topic),
			zap.Int32s("partitions", assignedPartitions),
		)
	}
}

func (cs consumer) onRevoked(
	_ context.Context,
	_ *kgo.Client,
	topicPartitions map[string][]int32,
) {
	for topic, revokedPartitions := range topicPartitions {
		if topic != cs.topic {
			cs.logger.Warn(
				"can't revoke partitions: not my topic",
				zap.String("got.topic", topic),
				zap.String("expected.topic", cs.topic),
			)

			continue
		}

		// TODO: commit offsets for this partition

		// It is VERY important to deregister the whole batch atomically.
		cs.partitionDispatcher.Stop(revokedPartitions...)

		cs.logger.Info(
			"kafka partitions revoked",
			zap.String("topic", topic),
			zap.Int32s("partitions", revokedPartitions),
		)
	}
}
