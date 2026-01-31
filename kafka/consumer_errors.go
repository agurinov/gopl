package kafka

import (
	"cmp"
	"context"
	"errors"

	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"

	"github.com/agurinov/gopl/graceful"
)

type eventloopAction uint8

// NOTE:
// skipAction means that error is retryable.
// exitAction means that error is fatal and eventloop should be interrupted.
// waitAction means that partition is IDLE and we should not poll again in exponential manner.
//
// TODO: implement more robust condition about IDLE
const (
	processAction eventloopAction = iota
	skipAction
	exitAction
	waitAction
)

func (consumer) analyzeFetches(
	fetches kgo.Fetches,
	l *zap.Logger,
) eventloopAction {
	isClosed := cmp.Or(
		fetches.IsClientClosed(),
		graceful.IsClosed(),
	)
	if isClosed {
		return exitAction
	}

	err := fetches.Err0()

	var (
		dataLossErr     *kgo.ErrDataLoss
		groupSessionErr *kgo.ErrGroupSession
		isIdle          = errors.Is(err, context.DeadlineExceeded)
	)

	switch {
	case err == nil:
		return processAction
	case errors.Is(err, context.Canceled):
		return exitAction
	case isIdle:
		return waitAction
	case errors.As(err, &dataLossErr):
		// Consumer will recover, the log is needed for further investigation for why
		// data loss had occurred.
		l.Warn(
			"data loss occurred, thus Kafka consumer was reset",
			zap.Error(dataLossErr),
		)

		return skipAction
	case errors.As(err, &groupSessionErr):
		l.Warn(
			"consumer group member was kicked from or was never able to join the group",
			zap.Error(groupSessionErr),
		)

		return skipAction
	default:
		l.Error(
			"unexpected fetch error",
			zap.Error(err),
		)

		return exitAction
	}
}
