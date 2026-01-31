package kafka

type (
	consumerMetrics struct {
		// idle backoff wait
		// client switch lock
		// poll partition kgo
		// handler work
		// discard work
		// rebalance
	}
)

func newConsumerMetrics() consumerMetrics {
	return consumerMetrics{}
}
