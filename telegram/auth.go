package telegram

import (
	"context"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	initdata "github.com/telegram-mini-apps/init-data-golang"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	c "github.com/agurinov/gopl/patterns/creational"
)

type (
	Auth struct {
		logger       *zap.Logger
		botTokens    map[string]string
		dummyEnabled bool
	}
	AuthOption c.Option[Auth]
)

const (
	tmaAuthSchema   = "tma"
	authHeaderParts = 2
)

var NewAuth = c.NewWithValidate[Auth, AuthOption]

func (a Auth) authFunc(initDataString string) (User, error) {
	if a.dummyEnabled {
		a.logger.Warn(
			"authenticated user",
			zap.Bool("dummy", true),
		)

		return Dummy(), nil
	}

	var validateErr error

LOOP:
	for botName, botToken := range a.botTokens {
		validateErr = initdata.Validate(initDataString, botToken, 0)
		switch validateErr {
		case nil:
			a.logger.Debug(
				"authenticated user",
				zap.String("bot_name", botName),
			)

			break LOOP
		default:
			a.logger.Debug(
				"can't authenticate user",
				zap.String("bot_name", botName),
				zap.Error(validateErr),
			)
		}
	}

	if validateErr != nil {
		return User{}, validateErr
	}

	initData, err := initdata.Parse(initDataString)
	if err != nil {
		return User{}, err
	}

	return initData.User, nil
}

func (a Auth) UnaryServerInterceptor(
	ctx context.Context,
	in any,
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	initDataString, err := auth.AuthFromMD(ctx, tmaAuthSchema)
	if err != nil {
		return ctx, err
	}

	user, err := a.authFunc(initDataString)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	ctx = SetUser(ctx, user)

	return handler(ctx, in)
}

func (a Auth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var initDataString string

		switch parts := strings.SplitN(
			r.Header.Get("Authorization"),
			" ",
			authHeaderParts,
		); {
		case len(parts) != authHeaderParts:
			http.Error(w, "", http.StatusUnauthorized)

			return
		case parts[0] != tmaAuthSchema:
			http.Error(w, "", http.StatusUnauthorized)

			return
		default:
			initDataString = parts[1]
		}

		user, err := a.authFunc(initDataString)
		if err != nil {
			http.Error(w, "", http.StatusUnauthorized)

			return
		}

		ctx := r.Context()
		ctx = SetUser(ctx, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
