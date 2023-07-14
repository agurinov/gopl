package kafka

import (
	"context"
	"fmt"

	"github.com/agurinov/gopl/backoff"
	pl_errors "github.com/agurinov/gopl/errors"
)

// healtcheck ensures that workflow is still reasonable and can be continued or retried
func healthcheck(ctx context.Context, b *backoff.Backoff, err error) error {
	if err != nil && !pl_errors.IsRetryable(err) {
		return fmt.Errorf("permanent error: %w", err)
	}

	if b == nil {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			return nil
		}
	}

	backoffStat, backoffErr := b.Wait(ctx)
	if backoffErr != nil {
		return pl_errors.Or(err, backoffErr)
	}

	// TODO(a.gurinov): write metric
	_ = backoffStat

	return nil
}
