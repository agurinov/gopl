package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Init(_ string) error {
	var (
		oldGoCollector = collectors.NewGoCollector()
		newGoCollector = collectors.NewGoCollector(
			collectors.WithGoCollectorMemStatsMetricsDisabled(),
			collectors.WithGoCollectorRuntimeMetrics(collectors.MetricsAll),
		)
	)

	prometheus.Unregister(oldGoCollector)

	if err := prometheus.Register(newGoCollector); err != nil {
		return err
	}

	return nil
}

func Handler() http.Handler {
	return promhttp.Handler()
}
