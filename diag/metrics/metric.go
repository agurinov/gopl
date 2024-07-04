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
