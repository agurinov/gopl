package graceful

import (
	"context"

	"golang.org/x/sync/errgroup"
)

func runGroup(
	ctx context.Context,
	errCh chan error,
	stack []Closure,
) error {
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

	return g.Wait()
}
