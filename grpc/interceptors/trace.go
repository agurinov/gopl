package interceptors

import (
	"context"
	"strings"

	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/agurinov/gopl/diag/trace"
)

func TraceUnaryServerInterceptor(
	ctx context.Context,
	in any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (
	any,
	error,
) {
	ctx, span := trace.StartSpan(ctx, "grpc.service")
	defer span.End()

	serviceName, methodName, _ := strings.Cut(
		strings.TrimPrefix(info.FullMethod, "/"),
		"/",
	)

	out, err := handler(ctx, in)

	span.SetAttributes(
		semconv.RPCSystemGRPC,
		semconv.RPCService(serviceName),
		semconv.RPCMethod(methodName),
		semconv.RPCGRPCStatusCodeKey.Int(
			int(status.Code(err)),
		),
	)

	if err != nil {
		return nil, trace.CatchError(span, err)
	}

	return out, nil
}
