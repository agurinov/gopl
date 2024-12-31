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
	basic struct {
		handler           http.Handler
		logger            *zap.Logger
		customMiddlewares []middlewares.Middleware
	}
	BasicOption c.Option[basic]
)

var NewBasic = c.NewWithValidate[basic, BasicOption]

func (h basic) Handler() http.Handler {
	r := chi.NewRouter()

	r.Use(
		middlewares.Trace,
		middlewares.Metrics(
			metrics.WithBuckets(metrics.BucketFast),
		),
		middlewares.AccessLog(
			h.logger,
			zapcore.InfoLevel,
		),
	)

	r.Use(h.customMiddlewares...)

	r.Handle("/*", h.handler)

	return r
}
