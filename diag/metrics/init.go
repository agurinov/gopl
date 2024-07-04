package metrics

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

var (
	cmdName    string
	registerer = prometheus.DefaultRegisterer
)

func Init(cn string) error {
	cmdName = cn

	var (
		oldGoCollector = collectors.NewGoCollector()
		newGoCollector = collectors.NewGoCollector(
			collectors.WithGoCollectorMemStatsMetricsDisabled(),
			collectors.WithGoCollectorRuntimeMetrics(collectors.MetricsAll),
		)
	)

	registerer.Unregister(oldGoCollector)

	if err := registerer.Register(newGoCollector); err != nil {
		return fmt.Errorf(
			"can't init metrics: %w",
			err,
		)
	}

	registerer = prometheus.WrapRegistererWith(
		prometheus.Labels{"service": cmdName},
		registerer,
	)

	return nil
}
