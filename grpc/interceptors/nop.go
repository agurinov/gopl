package interceptors

import (
	"context"

	"google.golang.org/grpc"
)

func NopUnaryServerInterceptor(
	ctx context.Context,
	in any,
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	return handler(ctx, in)
}
