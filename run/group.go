package run

import (
	"context"
	"errors"

	"golang.org/x/sync/errgroup"

	"github.com/agurinov/gopl/x"
)

func GroupSoft(
	ctx context.Context,
	stack ...Fn,
) error {
	errCh := make(chan error, len(stack))
	defer close(errCh)

	g, gCtx := errgroup.WithContext(ctx)

	for _, f := range stack {
		if f == nil {
			continue
		}

		g.Go(func() error {
			errCh <- f(gCtx)

			return nil
		})
	}

	var (
		groupErr = g.Wait()
		chanErr  = x.FlattenErrors(errCh)
	)

	return errors.Join(groupErr, chanErr)
}

func Group(
	ctx context.Context,
	stack ...Fn,
) error {
	g, gCtx := errgroup.WithContext(ctx)

	for _, f := range stack {
		if f == nil {
			continue
		}

		g.Go(func() error {
			return f(gCtx)
		})
	}

	return g.Wait()
}
