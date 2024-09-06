package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/agurinov/gopl/diag/metrics"
	"github.com/agurinov/gopl/diag/probes"
	"github.com/agurinov/gopl/http/middlewares"
	c "github.com/agurinov/gopl/patterns/creational"
)

type (
	debug struct {
		atomicLogLevel *zap.AtomicLevel
		logger         *zap.Logger
		prober         *probes.Prober
	}
	DebugOption c.Option[debug]
)

var NewDebug = c.NewWithValidate[debug, DebugOption]

func (h debug) Handler() http.Handler {
	r := chi.NewRouter()

	r.Use(
		middlewares.Trace,
		middlewares.Metrics(
			metrics.WithBuckets(metrics.BucketFast),
		),
		middlewares.AccessLog(h.logger),
		// middlewares.Panic(obj.logger),
	)

	if h.atomicLogLevel != nil {
		r.Mount("/logger", *h.atomicLogLevel)
	}

	r.Mount("/debug", middleware.Profiler())
	r.Mount("/metrics", Metrics())

	if h.prober != nil {
		r.Get("/probes/startup", probeHandler(h.prober.Startup))
		r.Get("/probes/readiness", probeHandler(h.prober.Readiness))
		r.Get("/probes/liveness", probeHandler(h.prober.Liveness))
	}

	return r
}

func probeHandler(probeGetter func() bool) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		switch probeGetter() {
		case true:
			w.WriteHeader(http.StatusOK)
		case false:
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	})
}
