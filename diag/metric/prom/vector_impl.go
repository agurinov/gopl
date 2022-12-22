package pl_prom

import (
	p "github.com/prometheus/client_golang/prometheus"
)

type vector interface {
	p.CounterVec | p.GaugeVec | p.HistogramVec
}

type vectorOption interface {
	p.CounterOpts | p.GaugeOpts | p.HistogramOpts
}

func newVector[V vector, O vectorOption](
	namespace string,
	component string,
	name string,
	help string,
	vectorFactory func(O, []string) *V,
	optionFactory func(O) O,
) (*V, error) {
	opts := O{
		Namespace: namespace,
		Subsystem: component,
		Name:      name,
		Help:      help,
	}

	if optionFactory != nil {
		opts = optionFactory(opts)
	}

	vec := vectorFactory(
		opts,
		[]string{},
	)

	// TODO(a.gurinov): register in registry

	return vec, nil
}
