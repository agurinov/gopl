package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/agurinov/gopl/diag/metrics"
	"github.com/agurinov/gopl/http/middlewares"
	c "github.com/agurinov/gopl/patterns/creational"
)

type (
	Basic struct {
		handlers          map[string]http.Handler
		logger            *zap.Logger
		customMiddlewares []middlewares.Middleware
		handlersAutomount bool
		accessLogLevel    zapcore.Level
	}
	BasicOption c.Option[Basic]
)

func NewBasic(opts ...BasicOption) (Basic, error) {
	h := Basic{
		accessLogLevel:    zapcore.InfoLevel,
		handlersAutomount: true,
	}

	return c.ConstructWithValidate(h, opts...)
}

func (h Basic) Handler() http.Handler { return h.Router() }

func (h Basic) Router() chi.Router {
	r := chi.NewRouter()

	r.Use(
		middlewares.Trace,
		middlewares.Metrics(
			metrics.WithBuckets(metrics.BucketFast),
		),
		middlewares.AccessLog(
			h.logger,
			h.accessLogLevel,
		),
	)

	r.Use(h.customMiddlewares...)

	for path, handler := range h.handlers {
		if !h.handlersAutomount {
			continue
		}

		r.Handle(path, handler)

		h.logger.Info(
			"handler mounted",
			zap.String("path", path),
		)
	}

	return r
}
