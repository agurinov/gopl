package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/agurinov/gopl/http/middlewares"
	c "github.com/agurinov/gopl/patterns/creational"
)

type (
	debug struct {
		atomicLogLevel *zap.AtomicLevel
		logger         *zap.Logger
	}
	DebugOption c.Option[debug]
)

var NewDebug = c.NewWithValidate[debug, DebugOption]

func (h debug) Handler() http.Handler {
	r := chi.NewRouter()

	r.Use(
		// middlewares.Trace,
		// middlewares.Metrics(),
		middlewares.AccessLog(h.logger),
		// middlewares.Panic(obj.logger),
	)

	if h.atomicLogLevel != nil {
		r.Mount("/logger", *h.atomicLogLevel)
	}

	r.Mount("/debug", middleware.Profiler())
	r.Mount("/metrics", Metrics())
	// r.Mount("/health", handlers.Probes())

	return r
}
