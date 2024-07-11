package telegram_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/metadata"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/agurinov/gopl/telegram"
	pl_testing "github.com/agurinov/gopl/testing"
)

func TestAuth_Middleware(t *testing.T) {
	type (
		args struct {
			request      *http.Request
			dummyEnabled bool
		}
		results struct {
			statusCode int
			content    string
		}
	)

	newRequest := func(authHeader string) *http.Request {
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		request.Header.Set("Authorization", authHeader)

		return request
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := telegram.GetUser(r.Context())
		if err != nil {
			http.Error(w, "oops", http.StatusInternalServerError)
		}

		io.WriteString(w, user.Username) //nolint:errcheck
	})

	cases := map[string]struct {
		args    args
		results results
		pl_testing.TestCase
	}{
		"case00: no header": {
			args: args{
				request: newRequest(""),
			},
			results: results{
				statusCode: http.StatusUnauthorized,
				content:    "\n",
			},
		},
		"case01: wrong token": {
			args: args{
				request: newRequest("tma foobar"),
			},
			results: results{
				statusCode: http.StatusUnauthorized,
				content:    "\n",
			},
		},
		"case02: dummy: no header": {
			args: args{
				request:      newRequest(""),
				dummyEnabled: true,
			},
			results: results{
				statusCode: http.StatusUnauthorized,
				content:    "\n",
			},
		},
		"case02: dummy: wrong schema": {
			args: args{
				request:      newRequest("bearer foobar"),
				dummyEnabled: true,
			},
			results: results{
				statusCode: http.StatusUnauthorized,
				content:    "\n",
			},
		},
		"case03: dummy: wrong token": {
			args: args{
				request:      newRequest("tma foobar"),
				dummyEnabled: true,
			},
			results: results{
				statusCode: http.StatusOK,
				content:    telegram.Dummy().Username,
			},
		},
	}

	for name := range cases {
		name, tc := name, cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			auth, err := telegram.NewAuth(
				telegram.WithAuthLogger(zaptest.NewLogger(t)),
				telegram.WithAuthDummy(tc.args.dummyEnabled),
				telegram.WithAuthBotTokens(map[string]string{"FooBot": "foo"}),
			)
			require.NoError(t, err)
			require.NotNil(t, auth)

			var (
				recorder    = httptest.NewRecorder()
				authHandler = auth.Middleware(handler)
			)

			authHandler.ServeHTTP(recorder, tc.args.request)

			require.Equal(t, tc.results.statusCode, recorder.Code)
			require.Equal(t, tc.results.content, recorder.Body.String())
		})
	}
}

func TestAuth_Interceptor(t *testing.T) {
	type (
		args struct {
			ctx          context.Context
			dummyEnabled bool
		}
		results struct {
			statusCode codes.Code
			out        any
		}
	)

	newCtx := func(authHeader string) context.Context {
		md := metadata.MD{}
		md = md.Set("authorization", authHeader)

		return md.ToIncoming(context.Background())
	}

	handler := func(ctx context.Context, req any) (any, error) {
		user, err := telegram.GetUser(ctx)
		if err != nil {
			return nil, status.Error(codes.Internal, "oops")
		}

		return user.Username, nil
	}

	cases := map[string]struct {
		args    args
		results results
		pl_testing.TestCase
	}{
		"case00: no header": {
			args: args{
				ctx: newCtx(""),
			},
			results: results{
				statusCode: codes.Unauthenticated,
				out:        nil,
			},
		},
		"case01: wrong token": {
			args: args{
				ctx: newCtx("tma foobar"),
			},
			results: results{
				statusCode: codes.Unauthenticated,
				out:        nil,
			},
		},
		"case02: dummy: no header": {
			args: args{
				ctx:          newCtx(""),
				dummyEnabled: true,
			},
			results: results{
				statusCode: codes.Unauthenticated,
			},
		},
		"case02: dummy: wrong schema": {
			args: args{
				ctx:          newCtx("bearer foobar"),
				dummyEnabled: true,
			},
			results: results{
				statusCode: codes.Unauthenticated,
				out:        nil,
			},
		},
		"case03: dummy: wrong token": {
			args: args{
				ctx:          newCtx("tma foobar"),
				dummyEnabled: true,
			},
			results: results{
				statusCode: codes.OK,
				out:        telegram.Dummy().Username,
			},
		},
	}

	for name := range cases {
		name, tc := name, cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			auth, err := telegram.NewAuth(
				telegram.WithAuthLogger(zaptest.NewLogger(t)),
				telegram.WithAuthDummy(tc.args.dummyEnabled),
				telegram.WithAuthBotTokens(map[string]string{"FooBot": "foo"}),
			)
			require.NoError(t, err)
			require.NotNil(t, auth)

			out, err := auth.UnaryServerInterceptor(
				tc.args.ctx,
				nil,
				new(grpc.UnaryServerInfo),
				handler,
			)
			code := status.Code(err)

			require.Equal(t, tc.results.statusCode, code)
			require.Equal(t, tc.results.out, out)
		})
	}
}
