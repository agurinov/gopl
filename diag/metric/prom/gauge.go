package pl_prom

import (
	p "github.com/prometheus/client_golang/prometheus"
)

func NewGauge() (Metric[p.Gauge], error) {
	return newVector[p.GaugeVec, p.GaugeOpts](
		"namespace",
		"component",
		"metric_name",
		"some useful metric",
		p.NewGaugeVec,
		nil,
	)
}
