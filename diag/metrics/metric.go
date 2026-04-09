package metrics

import "github.com/prometheus/client_golang/prometheus"

func NewCounter(
	name string,
	labels []string,
	opts ...Option,
) *prometheus.CounterVec {
	return newCreator(opts...).newCounter(name, labels...)
}

func NewHistogram(
	name string,
	labels []string,
	opts ...Option,
) *prometheus.HistogramVec {
	return newCreator(opts...).newHistogram(name, labels...)
}

func NewGauge(
	name string,
	labels []string,
	opts ...Option,
) *prometheus.GaugeVec {
	return newCreator(opts...).newGauge(name, labels...)
}
