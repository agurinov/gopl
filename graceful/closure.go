package graceful

import "context"

type Closure func(ctx context.Context) error

func SimpleClosure(
	fn func(),
) Closure {
	return func(_ context.Context) error {
		fn()

		return nil
	}
}

func ErrorClosure(
	fn func() error,
) Closure {
	return func(_ context.Context) error {
		return fn()
	}
}
