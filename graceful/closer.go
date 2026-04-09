package graceful

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

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
	if len(cl.stack1)+len(cl.stack2) == 0 {
		cl.logger.Info("closer finished; no closures registered")

		return nil
	}

	// Wait OS or orchestrator to start closing gracefully all components.
	<-runCtx.Done()

	cl.logger.Info(
		"closer started; going to run closures",
		zap.Stringer("timeout", cl.timeout),
		zap.Dict(
			"1st wave",
			zap.Int("closures", len(cl.stack1)),
		),
		zap.Dict(
			"2nd wave",
			zap.Int("closures", len(cl.stack2)),
		),
	)

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

	err := errors.Join(g1, g2)

	lvl := zapcore.InfoLevel
	if err != nil {
		lvl = zapcore.ErrorLevel
	}

	cl.logger.Log(lvl,
		"closer finished",
		zap.Error(err),
	)

	return err
}
