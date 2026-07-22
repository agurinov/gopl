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

func (cr creator) register(vec prometheus.Collector) prometheus.Collector {
	if !cr.useExisting {
		registerer.MustRegister(vec)

		return vec
	}

	var (
		existsErr = new(prometheus.AlreadyRegisteredError)
		err       = registerer.Register(vec)
	)

	switch {
	case errors.As(err, existsErr):
		return existsErr.ExistingCollector
	case err != nil:
		panic(err)
	}

	return vec
}

func (cr creator) newCounter(name string, labels ...string) *prometheus.CounterVec {
	metricName := cr.metricName(name)

	vec := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: metricName,
		},
		labels,
	)

	return registerAs[*prometheus.CounterVec](cr, vec)
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

	return registerAs[*prometheus.HistogramVec](cr, vec)
}

func (cr creator) newGauge(name string, labels ...string) *prometheus.GaugeVec {
	metricName := cr.metricName(name)

	vec := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: metricName,
		},
		labels,
	)

	return registerAs[*prometheus.GaugeVec](cr, vec)
}
