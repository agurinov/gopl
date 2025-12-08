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
	return func(graceCtx context.Context) error {
		if !w.closed.CompareAndSwap(false, true) {
			w.logger.Info("graceful stop already called; skipping")

			return nil
		}

		<-graceCtx.Done()

		ctx := w.ctx

		if err := f(ctx); err != nil {
			return err
		}

		w.forceCancel()

		return nil
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

func (w Wrapper) WrapRun(f Closure) Closure {
	return func(context.Context) error {
		ctx := w.ctx

		for {
			if w.IsClosed(ctx) {
				w.logger.Info("gracefully stop Run")

				return nil
			}

			//nolint:contextcheck
			if err := f(ctx); err != nil {
				return err
			}
		}
	}
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
