package telegram

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/agurinov/gopl/diag/trace"
)

type ctxKey string

var userCtxKey ctxKey = "user"

func SetUser(ctx context.Context, user User) context.Context {
	return context.WithValue(ctx, userCtxKey, user)
}

func GetUser(ctx context.Context) (User, error) {
	ctx, span := trace.StartSpan(ctx, "telegram.auth")
	defer span.End()

	user, ok := ctx.Value(userCtxKey).(User)
	if !ok {
		return User{}, trace.CatchError(span,
			status.Errorf(codes.Unauthenticated, "context is not authenticated"),
		)
	}

	span.SetAttributes(
		attribute.Int64("enduser.tg_id", user.ID),
	)

	return user, nil
}
