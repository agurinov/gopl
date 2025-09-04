package graceful

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	c "github.com/agurinov/gopl/patterns/creational"
)

type (
	closeF func(ctx context.Context) error
	Closer struct {
		logger  *zap.Logger
		stack   []closeF
		timeout time.Duration
	}
	CloserOption c.Option[Closer]
)

var NewCloser = c.NewWithValidate[Closer, CloserOption]

func (cl *Closer) AddCloser(
	fn func(),
) {
	if fn == nil {
		return
	}

	cl.stack = append(cl.stack, func(_ context.Context) error {
		fn()

		return nil
	})
}

func (cl *Closer) AddErrorCloser(
	fn func() error,
) {
	if fn == nil {
		return
	}

	cl.stack = append(cl.stack, func(_ context.Context) error {
		return fn()
	})
}

func (cl *Closer) AddContextErrorCloser(
	fn func(context.Context) error,
) {
	if fn == nil {
		return
	}

	cl.stack = append(cl.stack, fn)
}

//nolint:contextcheck
func (cl *Closer) WaitForShutdown(ctx context.Context) error {
	<-ctx.Done()

	cl.logger.Info(
		"shutting down closer functions",
		zap.Int("count", len(cl.stack)),
		zap.Stringer("timeout", cl.timeout),
	)

	if len(cl.stack) == 0 {
		return nil
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(
		context.Background(),
		cl.timeout,
	)
	defer shutdownCancel()

	errCh := make(chan error, len(cl.stack))
	defer close(errCh)

	g, gCtx := errgroup.WithContext(shutdownCtx)

	for _, f := range cl.stack {
		g.Go(func() error {
			errCh <- f(gCtx)

			return nil
		})
	}

	if waitErr := g.Wait(); waitErr != nil {
		return waitErr
	}

	errs := make([]error, 0, len(cl.stack))

	for err := range errCh {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}
