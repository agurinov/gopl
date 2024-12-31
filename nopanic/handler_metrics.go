package nopanic

import "github.com/prometheus/client_golang/prometheus"

type handlerMetrics struct {
	panicRecovered *prometheus.CounterVec
}

func (m handlerMetrics) recoveredPanicInc() {
	if m.panicRecovered == nil {
		return
	}
}
