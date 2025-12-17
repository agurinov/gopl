package graceful

import (
	"context"
	"sync/atomic"

	"go.uber.org/zap"

	c "github.com/agurinov/gopl/patterns/creational"
)

type (
	Wrapper struct {
		logger      *zap.Logger
		ctx         context.Context //nolint:containedctx
		forceCancel context.CancelFunc
		closed      *atomic.Bool
	}
	WrapperOption c.Option[Wrapper]
)

func (w Wrapper) WrapClose(f Closure) Closure {
	l := w.logger.With(
		zap.Stringer("closure", f),
	)

	l.Info("wrapping Close")

	return func(graceCtx context.Context) error {
		if isFirstClose := w.closed.CompareAndSwap(false, true); isFirstClose {
			w.logger.Info("grace cordon enabled")
		}

		<-graceCtx.Done()

		l.Info("grace period passed; force exiting")

		//nolint:contextcheck
		if err := f(w.ctx); err != nil {
			return err
		}

		w.forceCancel()

		return nil
	}
}

func (w Wrapper) WrapRun(f Closure) Closure {
	l := w.logger.With(
		zap.Stringer("closure", f),
	)

	l.Info("wrapping Run")

	return func(context.Context) error {
		ctx := w.ctx

		for {
			if w.closed.Load() {
				l.Info("gracefully stopped Run")

				return nil
			}

			select {
			case <-w.ctx.Done():
				l.Warn("force stopped Run")

				return nil
			default:
			}

			//nolint:contextcheck
			if err := f(ctx); err != nil {
				return err
			}
		}
	}
}

func (w Wrapper) IsClosed(ctxs ...context.Context) bool {
	if w.closed.Load() {
		return true
	}

	ctxs = append(ctxs, w.ctx)

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
	ctx, cancel := context.WithCancel(context.Background())

	obj := Wrapper{
		ctx:         ctx,
		forceCancel: cancel,
		closed:      new(atomic.Bool),
	}

	return c.ConstructWithValidate(obj, opts...)
}
