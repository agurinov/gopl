package pl_prom

import (
	p "github.com/prometheus/client_golang/prometheus"
)

func NewHistogram() (Metric[p.Observer], error) {
	return newVector[p.HistogramVec, p.HistogramOpts](
		"namespace",
		"component",
		"metric_name",
		"some useful metric",
		p.NewHistogramVec,
		func(o p.HistogramOpts) p.HistogramOpts {
			o.Buckets = p.DefBuckets

			return o
		},
	)
}
