package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	c "github.com/agurinov/gopl/patterns/creational"
)

type (
	debug struct {
		atomicLogLevel *zap.AtomicLevel
	}
	DebugOption c.Option[debug]
)

var NewDebug = c.New[debug, DebugOption]

func (h debug) Handler() http.Handler {
	r := chi.NewRouter()

	if h.atomicLogLevel != nil {
		r.Mount("/logger", *h.atomicLogLevel)
	}

	r.Mount("/debug", middleware.Profiler())

	// r.Mount("/metrics", handlers.Prometheus)
	// r.Mount("/health", handlers.Probes())

	return r
}
