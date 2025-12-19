package appcmd

import (
	"cmp"
	"context"
	"errors"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/agurinov/gopl/run"
)

func Start(
	ctx context.Context,
	logger *zap.Logger,
	stack ...run.Closure,
) {
	logger.Info("starting application")

	waitErr := run.Group(ctx, stack...)

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

// Deprecated: Use Start instead.
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
