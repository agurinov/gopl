package handlers

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/agurinov/gopl/http/middlewares"
)

func WithBasicLogger(logger *zap.Logger) BasicOption {
	return func(b *basic) error {
		if logger == nil {
			return nil
		}

		b.logger = logger.Named("http.handler.basic")

		return nil
	}
}

func WithBasicCustomMiddlewares(mw ...middlewares.Middleware) BasicOption {
	return func(b *basic) error {
		b.customMiddlewares = append(b.customMiddlewares, mw...)

		return nil
	}
}

func WithBasicHandler(h http.Handler) BasicOption {
	return func(b *basic) error {
		b.handler = h

		return nil
	}
}
