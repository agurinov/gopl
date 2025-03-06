package handlers

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/agurinov/gopl/http/middlewares"
)

func WithBasicLogger(logger *zap.Logger) BasicOption {
	return func(b *Basic) error {
		if logger == nil {
			return nil
		}

		b.logger = logger.Named("http.handler.basic")

		return nil
	}
}

func WithBasicCustomMiddlewares(mw ...middlewares.Middleware) BasicOption {
	return func(b *Basic) error {
		b.customMiddlewares = append(b.customMiddlewares, mw...)

		return nil
	}
}

func WithBasicHandlers(h map[string]http.Handler) BasicOption {
	return func(b *Basic) error {
		b.handlers = h

		return nil
	}
}
