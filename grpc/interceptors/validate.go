package interceptors

import (
	"context"

	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func ValidatorUnaryServer(v protovalidate.Validator) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		in any,
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		if protoMessage, ok := in.(proto.Message); ok {
			if err := v.Validate(protoMessage); err != nil {
				return nil, status.Error(
					codes.InvalidArgument,
					err.Error(),
				)
			}
		}

		return handler(ctx, in)
	}
}

func ValidatorUnaryClient(v protovalidate.Validator) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		in any,
		reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		if protoMessage, ok := in.(proto.Message); ok {
			if err := v.Validate(protoMessage); err != nil {
				return status.Error(
					codes.InvalidArgument,
					err.Error(),
				)
			}
		}

		return invoker(ctx, method, in, reply, cc, opts...)
	}
}
