package graceful

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	c "github.com/agurinov/gopl/patterns/creational"
)

type (
	closeF func(ctx context.Context) error
	Closer struct {
		logger  *zap.Logger
		stack1  []closeF
		stack2  []closeF
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
	fn func(),
	opts ...AddOption,
) {
	if fn == nil {
		return
	}

	closure := func(_ context.Context) error {
		fn()

		return nil
	}

	cl.AddContextErrorCloser(closure, opts...)
}

func (cl *Closer) AddErrorCloser(
	fn func() error,
	opts ...AddOption,
) {
	if fn == nil {
		return
	}

	closure := func(_ context.Context) error {
		return fn()
	}

	cl.AddContextErrorCloser(closure, opts...)
}

func (cl *Closer) AddContextErrorCloser(
	fn func(context.Context) error,
	opts ...AddOption,
) {
	if fn == nil {
		return
	}

	args, err := c.New(opts...)
	cl.logger.Warn(
		"can't construct add args",
		zap.Error(err),
	)

	switch args.wave {
	case FirstWave:
		cl.stack1 = append(cl.stack1, fn)
	default:
		cl.stack2 = append(cl.stack2, fn)
	}
}

//nolint:contextcheck
func (cl *Closer) WaitForShutdown(ctx context.Context) error {
	<-ctx.Done()

	cl.logger.Info(
		"closer started; going to run functions",
		zap.Stringer("timeout", cl.timeout),
		zap.Dict(
			"1st wave",
			zap.Int("functions", len(cl.stack1)),
		),
		zap.Dict(
			"2nd wave",
			zap.Int("functions", len(cl.stack2)),
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

	errCh := make(chan error, allLen)
	defer close(errCh)

	if err := runGroup(
		shutdownCtx,
		errCh,
		cl.stack1,
	); err != nil {
		return fmt.Errorf("can't close 1st wave: %w", err)
	}

	if err := runGroup(
		shutdownCtx,
		errCh,
		cl.stack2,
	); err != nil {
		return fmt.Errorf("can't close 2nd wave: %w", err)
	}

	if err := joinErrors(errCh); err != nil {
		return fmt.Errorf("can't join errors: %w", err)
	}

	cl.logger.Info("closer finished")

	return nil
}
