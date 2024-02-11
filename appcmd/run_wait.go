package appcmd

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func RunWait(g *errgroup.Group, logger *zap.Logger) {
	switch waitErr := g.Wait(); {
	case waitErr == nil,
		errors.Is(waitErr, context.Canceled),
		errors.Is(waitErr, context.DeadlineExceeded):
		logger.Info("application stopped")
	default:
		logger.Fatal("application crashed", zap.Error(waitErr))
	}
}
