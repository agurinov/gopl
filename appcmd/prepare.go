package appcmd

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"

	"github.com/agurinov/gopl/diag/log"
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

	goMemLimit, err := envvars.GoMemLimit.Value()
	if err != nil {
		return ctx, stop, logger, fmt.Errorf(
			"can't parse %s: %w",
			envvars.GoMemLimit.String(),
			err,
		)
	}

	goMaxProcs, err := envvars.GoMaxProcs.Value()
	if err != nil {
		return ctx, stop, logger, fmt.Errorf(
			"can't parse %s: %w",
			envvars.GoMaxProcs.String(),
			err,
		)
	}

	logger.Info(
		"resources from env",
		zap.Int(envvars.GoMaxProcs.String(), goMaxProcs),
		zap.String(envvars.GoMemLimit.String(), goMemLimit),
	)

	if _, maxprocsErr := maxprocs.Set(maxprocs.Logger(logger.Sugar().Infof)); maxprocsErr != nil {
		return ctx, stop, logger, maxprocsErr
	}

	// TODO(a.gurinov): k8s memlimit

	return ctx, stop, logger, nil
}
