package grpc

import (
	"context"

	"google.golang.org/grpc"
)

type tokenAuth map[string]string

const (
	authorizationHeader = "authorization"
	bearerAuthSchema    = "Bearer" //nolint:gosec
)

func (t tokenAuth) GetRequestMetadata(
	context.Context,
	...string,
) (
	map[string]string,
	error,
) {
	return t, nil
}

func (tokenAuth) RequireTransportSecurity() bool { return false }

func AuthToken(token string) grpc.DialOption {
	return grpc.WithPerRPCCredentials(tokenAuth{
		authorizationHeader: bearerAuthSchema + " " + token,
	})
}
