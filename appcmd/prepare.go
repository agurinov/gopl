package appcmd

import (
	"context"
	"os/signal"
	"syscall"

	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"

	"github.com/agurinov/gopl/diag/log"
)

func Prepare(cmdName string) ( //nolint:revive
	context.Context,
	context.CancelFunc,
	*zap.Logger,
	error,
) {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	logger := log.MustNewZapSystem().Named("system").With(
		zap.String("cmd_name", cmdName),
	)

	if _, err := maxprocs.Set(maxprocs.Logger(logger.Sugar().Infof)); err != nil {
		return nil, nil, nil, err
	}

	// TODO(a.gurinov): k8s memlimit

	return ctx, stop, logger, nil
}
