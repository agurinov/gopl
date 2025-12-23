package graceful

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"

	c "github.com/agurinov/gopl/patterns/creational"
	"github.com/agurinov/gopl/run"
)

type (
	Closer struct {
		logger  *zap.Logger
		stack1  []run.Fn
		stack2  []run.Fn
		timeout time.Duration
	}
	CloserOption c.Option[Closer]
)

type (
	Wave    uint8
	addArgs struct {
		wave Wave
	}
	AddOption c.Option[addArgs]
)

const (
	SecondWave Wave = iota
	FirstWave
)

var NewCloser = c.NewWithValidate[Closer, CloserOption]

func (cl *Closer) AddCloser(
	closure run.Fn,
	opts ...AddOption,
) {
	if closure == nil {
		return
	}

	args, err := c.New(opts...)
	if err != nil {
		cl.logger.Warn(
			"can't construct add args",
			zap.Error(err),
		)
	}

	switch args.wave {
	case FirstWave:
		cl.stack1 = append(cl.stack1, closure)
	default:
		cl.stack2 = append(cl.stack2, closure)
	}
}

//nolint:contextcheck
func (cl *Closer) WaitForShutdown(runCtx context.Context) error {
	<-runCtx.Done()

	cl.logger.Info(
		"closer started; going to run closures",
		zap.Stringer("timeout", cl.timeout),
		zap.Dict(
			"1st wave",
			zap.Stringers("closures", cl.stack1),
		),
		zap.Dict(
			"2nd wave",
			zap.Stringers("closures", cl.stack2),
		),
	)

	allLen := len(cl.stack1) + len(cl.stack2)

	if allLen == 0 {
		return nil
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(
		context.Background(),
		cl.timeout,
	)
	defer shutdownCancel()

	g1 := run.GroupSoft(shutdownCtx, cl.stack1...)
	g2 := run.GroupSoft(shutdownCtx, cl.stack2...)

	if g1 != nil {
		g1 = fmt.Errorf("can't run 1st wave: %w", g1)
	}

	if g2 != nil {
		g2 = fmt.Errorf("can't run 2nd wave: %w", g2)
	}

	if err := errors.Join(g1, g2); err != nil {
		return err
	}

	cl.logger.Info("closer finished")

	return nil
}
