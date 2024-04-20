package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
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
