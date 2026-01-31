package metrics

import (
	"errors"
	"strings"

	"github.com/prometheus/client_golang/prometheus"

	c "github.com/agurinov/gopl/patterns/creational"
)

type (
	creator struct {
		buckets         []float64
		noServicePrefix bool
		useExisting     bool
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

func (cr creator) register(vec prometheus.Collector) {
	if !cr.useExisting {
		registerer.MustRegister(vec)

		return
	}

	var (
		existsErr = new(prometheus.AlreadyRegisteredError)
		err       = registerer.Register(vec)
	)

	switch {
	case errors.As(err, existsErr):
	case err != nil:
		panic(err)
	}
}

func (cr creator) newCounter(name string, labels ...string) *prometheus.CounterVec {
	metricName := cr.metricName(name)

	vec := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: metricName,
		},
		labels,
	)

	cr.register(vec)

	return vec
}

func (cr creator) newHistogram(name string, labels ...string) *prometheus.HistogramVec {
	metricName := cr.metricName(name)

	buckets := cr.buckets
	if len(buckets) == 0 {
		buckets = prometheus.DefBuckets
	}

	vec := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    metricName,
			Buckets: buckets,
		},
		labels,
	)

	cr.register(vec)

	return vec
}

func (cr creator) newGauge(name string, labels ...string) *prometheus.GaugeVec {
	metricName := cr.metricName(name)

	vec := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: metricName,
		},
		labels,
	)

	cr.register(vec)

	return vec
}
