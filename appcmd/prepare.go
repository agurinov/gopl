package appcmd

import (
	"context"
	"os/signal"
	"syscall"

	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"

	"github.com/agurinov/gopl/diag/log"
	"github.com/agurinov/gopl/diag/metrics"
	"github.com/agurinov/gopl/env/envvars"
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

	var (
		goMaxProcs, _ = envvars.GoMaxProcs.Value() //nolint:errcheck
		goMemLimit, _ = envvars.GoMemLimit.Value() //nolint:errcheck
	)

	logger.Info(
		"resources from env",
		zap.Int(envvars.GoMaxProcs.String(), goMaxProcs),
		zap.String(envvars.GoMemLimit.String(), goMemLimit),
	)

	if _, err := maxprocs.Set(maxprocs.Logger(logger.Sugar().Infof)); err != nil {
		return ctx, stop, logger, err
	}

	// TODO(a.gurinov): k8s memlimit

	if err := metrics.Init(cmdName); err != nil {
		return ctx, stop, logger, err
	}

	return ctx, stop, logger, nil
}
