package nopanic_test

import (
	"context"
	"sync"
	"testing"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/agurinov/gopl/nopanic"
	pl_testing "github.com/agurinov/gopl/testing"
)

func panicking() {
	panic("OOPS")
}

func TestHandler_UnaryServerInterceptor(t *testing.T) {
	pl_testing.Init(t)

	h, err := nopanic.NewHandler(
		nopanic.WithLogger(zaptest.NewLogger(t)),
	)
	require.NoError(t, err)
	require.NotNil(t, h)

	interceptor := h.UnaryServerInterceptor()
	require.NotNil(t, interceptor)

	type (
		args struct {
			grpcIn      any
			grpcHandler grpc.UnaryHandler
		}
		results struct {
			grpcStatusCode codes.Code
			grpcOut        any
		}
	)

	var (
		ctx            = context.TODO()
		successHandler = func(_ context.Context, _ any) (any, error) {
			return "out", nil
		}
		panicHandler = func(_ context.Context, _ any) (any, error) {
			panicking()

			return "out", nil
		}
	)

	cases := map[string]struct {
		args    args
		results results
		pl_testing.TestCase
	}{
		"case00: success handler": {
			args: args{
				grpcIn:      "in",
				grpcHandler: successHandler,
			},
			results: results{
				grpcStatusCode: codes.OK,
				grpcOut:        "out",
			},
		},
		"case01: panic recovered": {
			args: args{
				grpcIn:      "in",
				grpcHandler: panicHandler,
			},
			results: results{
				grpcStatusCode: codes.Internal,
			},
		},
	}

	for name := range cases {
		name, tc := name, cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			grpcOut, err := interceptor(
				ctx,
				"in",
				new(grpc.UnaryServerInfo),
				tc.args.grpcHandler,
			)

			grpcCode := status.Code(err)

			require.Equal(t, tc.results.grpcStatusCode, grpcCode)
			require.Equal(t, tc.results.grpcOut, grpcOut)
			require.Nil(t, recover())
		})
	}
}

func TestHandler_TelegramBotMiddleware(t *testing.T) {
	pl_testing.Init(t)

	h, err := nopanic.NewHandler(
		nopanic.WithLogger(zaptest.NewLogger(t)),
	)
	require.NoError(t, err)
	require.NotNil(t, h)

	middleware := h.TelegramBotMiddleware()
	require.NotNil(t, middleware)

	ctx := context.TODO()

	handler := middleware(
		func(_ context.Context, _ *bot.Bot, _ *models.Update) {
			panicking()
		},
	)
	handler(ctx, new(bot.Bot), new(models.Update))
	require.Nil(t, recover())
}

func TestHandler_Do(t *testing.T) {
	pl_testing.Init(t)

	h, err := nopanic.NewHandler(
		nopanic.WithLogger(zaptest.NewLogger(t)),
	)
	require.NoError(t, err)
	require.NotNil(t, h)

	h.Do(panicking)
	require.Nil(t, recover())
}

func TestHandler_Go(t *testing.T) {
	pl_testing.Init(t)

	h, err := nopanic.NewHandler(
		nopanic.WithLogger(zap.NewNop()),
	)
	require.NoError(t, err)
	require.NotNil(t, h)

	var wg sync.WaitGroup

	wg.Add(1)

	h.Go(func() {
		defer wg.Done()

		panicking()
	})

	wg.Wait()

	require.Nil(t, recover())
}
