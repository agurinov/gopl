package run

import (
	"context"
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/agurinov/gopl/run"
	"github.com/agurinov/gopl/x"
)

type (
	Dispatcher[K comparable] interface {
		Running() []K
		Run(context.Context, run.Fn, ...K)
		Stop(...K)
		GetContext(K) (context.Context, context.CancelCauseFunc)
	}
	//nolint:containedctx
	dispatcherContext struct {
		ctx    context.Context
		cancel context.CancelCauseFunc
	}
	dispatcher[K comparable] struct {
		group    *errgroup.Group
		contexts map[K]dispatcherContext
		mu       sync.RWMutex
	}
)

func (d *dispatcher[K]) Running() []K {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return x.MapKeys(d.contexts)
}

func (d *dispatcher[K]) stop(keys ...K) {
	for _, key := range keys {
		if cancel := d.contexts[key].cancel; cancel != nil {
			cancel(nil)
		}

		delete(d.contexts, key)
	}
}

func (d *dispatcher[K]) Stop(keys ...K) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.stop(keys...)
}

func (d *dispatcher[K]) GetContext(
	key K,
) (
	context.Context,
	context.CancelCauseFunc,
) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return d.contexts[key].ctx, d.contexts[key].cancel
}

func (d *dispatcher[K]) Run(
	ctx context.Context,
	runFn run.Fn,
	keys ...K,
) {
	if runFn == nil {
		return
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	d.stop(keys...)

	for _, key := range keys {
		tCtx, tCancel := context.WithCancelCause(ctx)

		d.group.Go(func() error {
			return runFn(tCtx)
		})

		d.contexts[key] = dispatcherContext{
			ctx:    tCtx,
			cancel: tCancel,
		}
	}
}

func NewDispatcher[T comparable]() (
	Dispatcher[T],
	error,
) {
	obj := dispatcher[T]{
		contexts: make(map[T]dispatcherContext),
		group:    new(errgroup.Group),
	}

	return &obj, nil
}
