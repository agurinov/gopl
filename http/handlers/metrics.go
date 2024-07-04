package handlers

import (
	"net/http"

	"github.com/agurinov/gopl/diag/metrics"
)

func Metrics() http.Handler {
	return metrics.Handler()
}
