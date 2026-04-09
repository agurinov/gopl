package run

import (
	"context"
	"fmt"
)

type contextKey string

func GetDispatcherKey[K comparable](
	ctx context.Context,
) (K, error) {
	var zero K

	key := contextKey(
		fmt.Sprintf("%T", zero),
	)

	value := ctx.Value(key)
	if value == nil {
		return zero, ErrNoContextKey
	}

	t, ok := value.(K)
	if !ok {
		return zero, ErrNoContextKey
	}

	return t, nil
}

func keyToContext[K comparable](
	ctx context.Context,
	k K,
) context.Context {
	key := contextKey(
		fmt.Sprintf("%T", k),
	)

	return context.WithValue(ctx, key, k)
}
