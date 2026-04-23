package kafka

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/agurinov/gopl/diag/metrics"
)

type (
	consumerMetrics struct {
		notMyRecord *prometheus.CounterVec
		idle        *prometheus.HistogramVec
		discarded   *prometheus.CounterVec
		committed   *prometheus.CounterVec
		polling     *prometheus.HistogramVec
		// handler work ??? via middleware

		// Separate operations??
		// client switch lock
		// rebalance
	}
)

func newConsumerMetrics() consumerMetrics {
	return consumerMetrics{
		notMyRecord: metrics.NewCounter(
			KafkaNotMyRecordCounterName,
			[]string{"topic", "partition"},
			metrics.WithoutServicePrefix(),
		),
		idle: metrics.NewHistogram(
			KafkaConsumerIdleDurationHistogramName,
			[]string{"topic", "partition"},
			metrics.WithoutServicePrefix(),
		),
		discarded: metrics.NewCounter(
			KafkaConsumerDiscardedCounterName,
			[]string{"topic", "partition"},
			metrics.WithoutServicePrefix(),
		),
		committed: metrics.NewCounter(
			KafkaConsumerCommittedCounterName,
			[]string{"topic", "partition"},
			metrics.WithoutServicePrefix(),
		),
		polling: metrics.NewHistogram(
			KafkaConsumerPollingDurationHistogramName,
			[]string{"topic", "partition"},
			metrics.WithoutServicePrefix(),
		),
	}
}

func (m consumerMetrics) discardedAdd(
	topic string,
	partition int32,
	cnt int,
) {
	if m.discarded == nil {
		return
	}

	m.discarded.WithLabelValues(
		topic,
		strconv.Itoa(int(partition)),
	).Add(
		float64(cnt),
	)
}

func (m consumerMetrics) committedAdd(
	topic string,
	partition int32,
	cnt int,
) {
	if m.committed == nil {
		return
	}

	m.committed.WithLabelValues(
		topic,
		strconv.Itoa(int(partition)),
	).Add(
		float64(cnt),
	)
}

func (m consumerMetrics) notMyRecordInc(
	topic string,
	partition int32,
) {
	if m.notMyRecord == nil {
		return
	}

	m.notMyRecord.WithLabelValues(
		topic,
		strconv.Itoa(int(partition)),
	).Inc()
}

func (m consumerMetrics) idleObserve(
	topic string,
	partition int32,
	elapsedTime time.Duration,
) {
	if m.idle == nil {
		return
	}

	m.idle.WithLabelValues(
		topic,
		strconv.Itoa(int(partition)),
	).Observe(
		elapsedTime.Seconds(),
	)
}

func (m consumerMetrics) pollingObserve(
	topic string,
	partition int32,
	elapsedTime time.Duration,
) {
	if m.polling == nil {
		return
	}

	m.polling.WithLabelValues(
		topic,
		strconv.Itoa(int(partition)),
	).Observe(
		elapsedTime.Seconds(),
	)
}
