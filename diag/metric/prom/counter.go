package pl_prom

import (
	p "github.com/prometheus/client_golang/prometheus"
)

func NewCounter() (Metric[p.Counter], error) {
	return newVector[p.CounterVec, p.CounterOpts](
		"namespace",
		"component",
		"metric_name",
		"some useful metric",
		p.NewCounterVec,
		nil,
	)
}
