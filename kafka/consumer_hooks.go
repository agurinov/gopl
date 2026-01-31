package kafka

import (
	"context"
	"slices"

	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"

	"github.com/agurinov/gopl/x"
)

func (cs consumer) OnFetchRecordBuffered(record Record) {
	key := string(record.Key)

	_ = key
	// TODO: detect multiple fetches for same record

	cs.logger.Warn(
		"record buffered multiple times",
		RecordLogFields(record)...,
	)
}

func (cs consumer) onAssigned(
	ctx context.Context,
	_ *kgo.Client,
	topicPartitions map[string][]int32,
) {
	for topic, assignedPartitions := range topicPartitions {
		l := cs.logger.With(
			zap.String("topic", topic),
		)

		if topic != cs.topic {
			l.Warn(
				"can't assign partitions: not my topic",
				zap.String("config.topic", cs.topic),
				zap.String("got.topic", topic),
			)

			continue
		}

		// Filter out partitions
		// TODO: remove it
		partitions := x.SliceFilter(
			assignedPartitions,
			func(p int32) bool {
				isSkip := slices.Contains([]int32{1, 5, 4}, p)

				if isSkip {
					l.Info(
						"skip partition from processing",
						zap.Int32("partition", p),
					)
				}

				return !isSkip
			},
		)

		// TODO: check atomicity in batch case one by one
		for _, partition := range partitions {
			cs.partitionDispatcher.Run(
				ctx,
				cs.partitionLoop(partition),
				partition,
			)
		}

		l.Info(
			"Kafka partitions assigned",
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
		l := cs.logger.With(
			zap.String("topic", topic),
		)

		if topic != cs.topic {
			l.Warn(
				"can't revoke partitions: not my topic",
				zap.String("config.topic", cs.topic),
				zap.String("got.topic", topic),
			)

			continue
		}

		// TODO: commit offsets for this partition

		cs.partitionDispatcher.Stop(revokedPartitions...)

		l.Info(
			"Kafka partitions revoked",
			zap.Int32s("partitions", revokedPartitions),
		)
	}
}
