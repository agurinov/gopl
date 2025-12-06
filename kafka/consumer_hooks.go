package kafka

import (
	"context"

	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
)

func (c *consumer) OnFetchRecordBuffered(r *kgo.Record) {
}

func (c *consumer) onAssigned(
	ctx context.Context,
	_ *kgo.Client,
	topicPartitions map[string][]int32,
) {
	for topic, assignedPartitions := range topicPartitions {
		l := c.logger.With(
			zap.String("topic", topic),
		)

		if topic != c.config.topic {
			l.Warn("tried to assign partitions from topic, that we don't consume")

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

			go c.processPartition(pCtx, partition)
		}

		l.Info(
			"Kafka partitions assigned",
			zap.Int32s("partitions", assignedPartitions),
		)
	}
}

func (c *consumer) onRevoked(
	_ context.Context,
	_ *kgo.Client,
	topicPartitions map[string][]int32,
) {
	for topic, revokedPartitions := range topicPartitions {
		l := c.logger.With(
			zap.String("topic", topic),
		)

		if topic != c.config.topic {
			l.Warn("revoked partitions from topic, that we don't consume")

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
