package appcmd

import (
	"cmp"
	"context"
	"errors"

	"github.com/agurinov/gopl/graceful"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func Start(
	ctx context.Context,
	logger *zap.Logger,
	stack ...graceful.Closure,
) {
	g, gCtx := errgroup.WithContext(ctx)

	for _, f := range stack {
		if f == nil {
			continue
		}

		g.Go(func() error {
			return f(gCtx)
		})
	}

	RunWait(g, logger)
}

func RunWait(g *errgroup.Group, logger *zap.Logger) {
	logger.Info("starting application")

	waitErr := g.Wait()

	isSuccess := cmp.Or(
		waitErr == nil,
		errors.Is(waitErr, context.Canceled),
		errors.Is(waitErr, context.DeadlineExceeded),
	)

	switch {
	case isSuccess:
		logger.Info("application stopped")
	default:
		logger.Fatal(
			"application crashed",
			zap.Error(waitErr),
		)
	}
}
