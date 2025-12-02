package graceful

import (
	"context"
	"errors"

	"golang.org/x/sync/errgroup"
)

func runGroup(
	ctx context.Context,
	errCh chan error,
	stack []closeF,
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

func joinErrors(
	errCh chan error,
) error {
	var joinedErr error

	for range cap(errCh) {
		select {
		case err := <-errCh:
			joinedErr = errors.Join(joinedErr, err)
		default:
		}
	}

	return joinedErr
}
