package telegram

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
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
		noBot        bool
	}
	AuthOption c.Option[Auth]
)

const tmaAuthSchema = "tma"

var NewAuth = c.NewWithValidate[Auth, AuthOption]

func (a Auth) authFunc(initDataString string) (User, error) {
	if a.dummyEnabled {
		a.logger.Warn(
			"authenticated user",
			zap.Bool("dummy", true),
		)

		return Dummy(), nil
	}

	var (
		validateErr  error
		authorityBot string
	)

LOOP:
	for botName, botToken := range a.botTokens {
		validateErr = initdata.Validate(initDataString, botToken, 0)
		switch validateErr {
		case nil:
			a.logger.Debug(
				"authenticated user",
				zap.String("bot_name", botName),
			)
			authorityBot = botName

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

	user := User{
		ID:           initData.User.ID,
		Username:     initData.User.Username,
		FirstName:    initData.User.FirstName,
		LastName:     initData.User.LastName,
		IsBot:        initData.User.IsBot,
		AuthorityBot: authorityBot,
	}

	if vErr := validator.New().Struct(user); vErr != nil {
		return User{}, vErr
	}

	if a.noBot && user.IsBot {
		return User{}, errors.New("can't authenticate bot")
	}

	return user, nil
}

func (a Auth) UnaryServerInterceptor(
	ctx context.Context,
	in any,
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	initDataString, err := auth.AuthFromMD(ctx, tmaAuthSchema)
	if err != nil && !a.dummyEnabled {
		return nil, err
	}

	user, err := a.authFunc(initDataString)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	ctx = SetUser(ctx, user)

	return handler(ctx, in)
}

func (Auth) authFromHeader(r *http.Request, expectedScheme string) (string, error) {
	val := r.Header.Get("Authorization")
	if val == "" {
		return "", errors.New("request unauthenticated with " + expectedScheme)
	}

	scheme, token, found := strings.Cut(val, " ")
	if !found {
		return "", errors.New("bad authorization string")
	}

	if !strings.EqualFold(scheme, expectedScheme) {
		return "", errors.New("request unauthenticated with " + expectedScheme)
	}

	return token, nil
}

func (a Auth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		initDataString, err := a.authFromHeader(r, tmaAuthSchema)
		if err != nil && !a.dummyEnabled {
			http.Error(w, "", http.StatusUnauthorized)

			return
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
