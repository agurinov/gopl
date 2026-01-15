package kafka

import (
	"context"

	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
)

func (c consumer[R, V]) OnFetchRecordBuffered(r *kgo.Record) {
}

func (c consumer[R, V]) onAssigned(
	ctx context.Context,
	_ *kgo.Client,
	topicPartitions map[string][]int32,
) {
	for topic, assignedPartitions := range topicPartitions {
		l := c.logger.With(
			zap.String("topic", topic),
		)

		if topic != c.config.topic {
			l.Warn(
				"can't assign partitions: not my topic",
				zap.String("config.topic", c.config.topic),
				zap.String("got.topic", topic),
			)

			continue
		}

		c.partitionHolder.assignPartitions(ctx, assignedPartitions)

		for _, partition := range assignedPartitions {
			switch partition {
			case 1, 5, 4:
			default:
				l.Info(
					"skip partition from processing",
					zap.Int32("partition", partition),
				)

				continue
			}

			pCtx := c.partitionHolder.partitionContext(partition)

			go func() {
				if err := c.processPartition(pCtx, partition); err != nil {
					l.Error(
						"can't process partition",
						zap.Int32("partition", partition),
						zap.Error(err),
					)
				}
			}()
		}

		l.Info(
			"Kafka partitions assigned",
			zap.Int32s("partitions", assignedPartitions),
		)
	}
}

func (c consumer[R, V]) onRevoked(
	_ context.Context,
	_ *kgo.Client,
	topicPartitions map[string][]int32,
) {
	for topic, revokedPartitions := range topicPartitions {
		l := c.logger.With(
			zap.String("topic", topic),
		)

		if topic != c.config.topic {
			l.Warn(
				"can't revoke partitions: not my topic",
				zap.String("config.topic", c.config.topic),
				zap.String("got.topic", topic),
			)

			continue
		}

		for _, p := range revokedPartitions {
			// TODO: commit offsets for this partition
			c.partitionHolder.revokePartition(p)
		}

		l.Info(
			"Kafka partitions revoked",
			zap.Int32s("partitions", revokedPartitions),
		)
	}
}
