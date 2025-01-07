package interceptors

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/agurinov/gopl/strings"
)

func ErrorConverterUnaryServer(
	ctx context.Context,
	in any,
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (
	any,
	error,
) {
	out, err := handler(ctx, in)
	if err == nil {
		return out, nil
	}

	if grpcStatus, ok := status.FromError(err); ok {
		return nil, grpcStatus.Err()
	}

	var (
		grpcCode    codes.Code
		grpcMessage = strings.RedactedPlaceholder
		validateErr = new(validator.ValidationErrors)
	)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		grpcCode = codes.NotFound
	case errors.As(err, validateErr):
		grpcCode = codes.InvalidArgument
		grpcMessage = err.Error()
	default:
		grpcCode = codes.Unknown
	}

	return nil, status.Error(grpcCode, grpcMessage)
}
