package middlewares

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/agurinov/gopl/diag/metrics"
)

var (
	onceServer         sync.Once
	histServerDuration *prometheus.HistogramVec
)

func Metrics(options ...metrics.Option) Middleware {
	onceServer.Do(func() {
		options = append(
			options,
			metrics.WithoutServicePrefix(),
		)
		histServerDuration = metrics.NewHistogram(
			metrics.HTTPServerDurationHistogramName,
			[]string{"method", "path", "status"},
			options...,
		)
	})

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			recorder := &statusRecorder{
				ResponseWriter: w,
				Status:         http.StatusOK,
			}

			startTime := time.Now()

			next.ServeHTTP(recorder, r)

			elapsedTime := time.Since(startTime)

			var (
				method = r.Method
				status = strconv.Itoa(recorder.Status)
				path   string
			)

			switch recorder.Status {
			case http.StatusNotFound:
				path = "not_found"
			default:
				path = "implement_me"
			}

			histServerDuration.WithLabelValues(
				method,
				path,
				status,
			).Observe(
				elapsedTime.Seconds(),
			)
		})
	}
}
