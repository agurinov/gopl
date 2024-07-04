package metrics

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus"

	c "github.com/agurinov/gopl/patterns/creational"
)

type (
	creator struct {
		buckets         []float64
		noServicePrefix bool
	}
	Option = c.Option[creator]
)

var newCreator = c.MustNew[creator, Option]

func (cr creator) metricName(name string) string {
	name = strings.ReplaceAll(name, "-", "_")

	if cr.noServicePrefix {
		return name
	}

	return strings.ReplaceAll(cmdName, "-", "_") + "_" + name
}

func (cr creator) newCounter(name string, labels ...string) *prometheus.CounterVec {
	metricName := cr.metricName(name)

	counterVec := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: metricName,
		},
		labels,
	)

	registerer.MustRegister(counterVec)

	return counterVec
}

func (cr creator) newHistogram(name string, labels ...string) *prometheus.HistogramVec {
	metricName := cr.metricName(name)

	if len(cr.buckets) == 0 {
		cr.buckets = prometheus.DefBuckets
	}

	histogramVec := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    metricName,
			Buckets: cr.buckets,
		},
		labels,
	)

	registerer.MustRegister(histogramVec)

	return histogramVec
}
