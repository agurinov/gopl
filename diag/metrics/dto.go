package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

func DTO(in prometheus.Metric) (*dto.Metric, error) {
	var dm dto.Metric

	if err := in.Write(&dm); err != nil {
		return nil, err
	}

	return &dm, nil
}
