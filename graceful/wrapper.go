package graceful

import (
	"context"
	"fmt"
	"io"
	"sync/atomic"

	"go.uber.org/zap"

	c "github.com/agurinov/gopl/patterns/creational"
)

type (
	runnable interface {
		Run(context.Context) error
	}
	Wrapper struct {
		logger      *zap.Logger
		ctx         context.Context //nolint:containedctx
		forceCancel context.CancelFunc
		closed      *atomic.Bool
	}
	WrapperOption c.Option[Wrapper]
)

func (w Wrapper) WrapClose(inner io.Closer) func(context.Context) error {
	return func(graceCtx context.Context) error {
		l := w.logger.With(
			zap.String("inner", fmt.Sprintf("%T", inner)),
		)

		if !w.closed.CompareAndSwap(false, true) {
			l.Info("graceful stop already called; skipping")

			return nil
		}

		<-graceCtx.Done()

		w.forceCancel()

		if err := inner.Close(); err != nil {
			return err
		}

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

func (w Wrapper) WrapRun(inner runnable) func(context.Context) error {
	return func(context.Context) error {
		var (
			ctx = w.ctx
			l   = w.logger.With(
				zap.String("inner", fmt.Sprintf("%T", inner)),
			)
		)

		for {
			if w.IsClosed(ctx) {
				l.Info("gracefully stop Run")

				return nil
			}

			//nolint:contextcheck
			if err := inner.Run(ctx); err != nil {
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
