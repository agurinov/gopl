package nopanic

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/agurinov/gopl/http/middlewares"
	c "github.com/agurinov/gopl/patterns/creational"
	"github.com/agurinov/gopl/telegram"
)

type (
	Handler struct {
		logger  *zap.Logger
		metrics handlerMetrics
	}
	Option c.Option[Handler]
)

var NewHandler = c.NewWithValidate[Handler, Option]

func (h Handler) safe(
	f func(),
	onPanic func(),
) {
	defer func() {
		if r := recover(); r != nil {
			var (
				err   = fmt.Errorf("%s", r)
				stack = string(debug.Stack())
			)

			h.logger.Error(
				"panic recovered",
				zap.String("stack", stack),
				zap.Error(err),
			)

			h.metrics.recoveredPanicInc()

			if onPanic != nil {
				onPanic()
			}
		}
	}()

	f()
}

func (h Handler) Middleware() middlewares.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.safe(
				func() {
					next.ServeHTTP(w, r)
				},
				func() {
					w.WriteHeader(http.StatusInternalServerError)
				},
			)
		})
	}
}

func (h Handler) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		in any,
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (
		out any,
		err error,
	) {
		h.safe(
			func() {
				out, err = handler(ctx, in)
			},
			func() {
				out, err = nil, status.Error(codes.Unknown, "")
			},
		)

		return out, err
	}
}

func (h Handler) TelegramBotMiddleware() telegram.BotMiddleware {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, bot *bot.Bot, update *models.Update) {
			h.safe(
				func() {
					next(ctx, bot, update)
				},
				nil,
			)
		}
	}
}

func (h *Handler) Go(f func()) {
	go func() {
		h.safe(f, nil)
	}()
}

func (h *Handler) Do(f func()) {
	h.safe(f, nil)
}
