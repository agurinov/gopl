package telegram

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ctxKey string

var userCtxKey ctxKey = "user"

func SetUser(ctx context.Context, user User) context.Context {
	return context.WithValue(ctx, userCtxKey, user)
}

func GetUser(ctx context.Context) (User, error) {
	user, ok := ctx.Value(userCtxKey).(User)
	if !ok {
		return User{}, status.Errorf(codes.Unauthenticated, "")
	}

	return user, nil
}
