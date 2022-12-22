package pl_prom

import (
	p "github.com/prometheus/client_golang/prometheus"
)

type metric interface {
	p.Counter
}

type Metric[M metric] interface {
	WithLabelValues(...string) M
}
