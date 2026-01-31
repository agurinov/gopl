package internal

import (
	"cmp"
	"context"
	"errors"
	"sync/atomic"

	"go.uber.org/zap"

	"github.com/agurinov/gopl/diag/log"
	c "github.com/agurinov/gopl/patterns/creational"
	"github.com/agurinov/gopl/run"
)

type (
	Wrapper interface {
		RunLoop(run.Fn) run.Fn
		Run(run.Fn) run.Fn
		Close(run.Fn) run.Fn
		IsClosed(...context.Context) bool
	}
	wrapper struct {
		logger      *zap.Logger
		safeCtx     context.Context //nolint:containedctx
		forceCancel context.CancelCauseFunc
		closed      *atomic.Bool
	}
	WrapperOption c.Option[wrapper]
)

func (w wrapper) Close(closeFn run.Fn) run.Fn {
	return func(graceCtx context.Context) error {
		if isFirstClose := w.closed.CompareAndSwap(false, true); isFirstClose {
			w.logger.Info("grace cordon enabled")
		}

		<-graceCtx.Done()

		w.logger.Info("grace period passed; force exiting")

		var closeErr error

		if closeFn != nil {
			// TODO: check this logic
			safeCtx := w.safeCtx

			//nolint:contextcheck
			closeErr = closeFn(safeCtx)
		}

		w.forceCancel(closeErr)

		return closeErr
	}
}

func (w wrapper) RunLoop(iterationFn run.Fn) run.Fn {
	return run.ErrorFn(func() error {
		safeCtx := w.safeCtx

		for {
			if w.closed.Load() {
				w.logger.Info("gracefully stopped Run")

				return nil
			}

			select {
			case <-safeCtx.Done():
				w.logger.Warn("force stopped Run")

				return nil
			default:
			}

			switch err := iterationFn(safeCtx); {
			case errors.Is(err, ErrStopLoop):
				return nil
			case err != nil:
				return err
			}
		}
	})
}

func (w wrapper) Run(runFn run.Fn) run.Fn {
	iterationFn := func(ctx context.Context) error {
		return cmp.Or(
			runFn(ctx),
			ErrStopLoop,
		)
	}

	return w.RunLoop(iterationFn)
}

func (w wrapper) IsClosed(ctxs ...context.Context) bool {
	if w.closed.Load() {
		return true
	}

	ctxs = append(ctxs, w.safeCtx)

	for i := range ctxs {
		if ctxs[i] == nil {
			continue
		}

		select {
		case <-ctxs[i].Done():
			return true
		default:
		}
	}

	return false
}

func NewWrapper(opts ...WrapperOption) (Wrapper, error) {
	ctx, cancel := context.WithCancelCause(context.Background())

	obj := wrapper{
		logger:      log.MustNewZapSystem(),
		safeCtx:     ctx,
		forceCancel: cancel,
		closed:      new(atomic.Bool),
	}

	return c.ConstructWithValidate(obj, opts...)
}
