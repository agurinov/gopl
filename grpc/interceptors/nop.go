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
) (
	any,
	error,
) {
	return handler(ctx, in)
}

func NopUnaryClientInterceptor(
	ctx context.Context,
	method string,
	in any,
	reply any,
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	return invoker(ctx, method, in, reply, cc, opts...)
}
