package nopanic

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/agurinov/gopl/diag/metrics"
)

type handlerMetrics struct {
	panicRecovered *prometheus.CounterVec
}

func newHandlerMetrics() handlerMetrics {
	return handlerMetrics{
		panicRecovered: metrics.NewCounter(
			metrics.NopanicHandlerCounterName,
			nil,
			metrics.WithoutServicePrefix(),
		),
	}
}

func (m handlerMetrics) recoveredPanicInc() {
	if m.panicRecovered == nil {
		return
	}
}
