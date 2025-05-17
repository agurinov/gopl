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
		logger           *zap.Logger
		botTokens        map[string]string
		noSignatureCheck bool
		noBotAllowed     bool
	}
	AuthOption c.Option[Auth]
)

const (
	tmaAuthSchema    = "tma"
	dummyBotUsername = "DummyBot"
)

var NewAuth = c.NewWithValidate[Auth, AuthOption]

func (a Auth) parseInitData(
	id initdata.InitData,
	authorityBot string,
) (User, error) {
	var privateChatID int64

	switch {
	case id.ChatType == initdata.ChatTypePrivate:
		privateChatID = id.ChatInstance
	case id.Chat.Type == initdata.ChatTypePrivate:
		privateChatID = id.Chat.ID
	case id.User.AllowsWriteToPm:
		privateChatID = id.User.ID
	}

	user := User{
		ID:           id.User.ID,
		Username:     id.User.Username,
		FirstName:    id.User.FirstName,
		LastName:     id.User.LastName,
		Photo:        id.User.PhotoURL,
		Language:     id.User.LanguageCode,
		IsBot:        id.User.IsBot,
		IsPremium:    id.User.IsPremium,
		AuthorityBot: authorityBot,
		PrivateChat: PrivateChat{
			ID:      privateChatID,
			Enabled: id.User.AllowsWriteToPm,
		},
	}

	if err := validator.New().Struct(user); err != nil {
		return User{}, err
	}

	if a.noBotAllowed && user.IsBot {
		return User{}, errors.New("can't authenticate bot")
	}

	a.logger.Debug(
		"authenticated user",
		zap.String("authority_bot", authorityBot),
		zap.String("tg_username", user.Username),
		zap.Int64("tg_id", user.ID),
	)

	return user, nil
}

func (a Auth) authFunc(initDataString string) (User, error) {
	initData, err := initdata.Parse(initDataString)
	if err != nil {
		return User{}, err
	}

	if a.noSignatureCheck {
		return a.parseInitData(initData, dummyBotUsername)
	}

	signatureErr := errors.New("no authority bots found")

	for botName, botToken := range a.botTokens {
		signatureErr = initdata.Validate(initDataString, botToken, 0)

		if signatureErr == nil {
			return a.parseInitData(initData, botName)
		}

		a.logger.Debug(
			"can't authenticate user, trying next",
			zap.String("authority_bot", botName),
			zap.String("tg_username", initData.User.Username),
			zap.Int64("tg_id", initData.User.ID),
			zap.Error(signatureErr),
		)
	}

	return User{}, signatureErr
}

func (a Auth) UnaryServerInterceptor(
	ctx context.Context,
	in any,
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	initDataString, err := auth.AuthFromMD(ctx, tmaAuthSchema)
	if err != nil {
		return nil, err
	}

	user, err := a.authFunc(initDataString)
	if err != nil {
		a.logger.Debug(
			"can't authenticate user",
			zap.Error(err),
		)

		return nil, status.Error(codes.Unauthenticated, "")
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
		if err != nil {
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
