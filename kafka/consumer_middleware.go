package kafka

import (
	"context"
)

type (
	Middleware  = func(Handler) Handler
	Middlewares []Middleware
)

func (mws Middlewares) Handler(h Handler) Handler {
	if len(mws) == 0 {
		return h
	}

	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}

	return h
}

func MetricsMiddleware(next Handler) Handler {
	return func(ctx context.Context, r Record) error {
		return next(ctx, r)
	}
}
