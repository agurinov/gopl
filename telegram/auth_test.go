package telegram_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/metadata"
	"github.com/stretchr/testify/require"
	initdata "github.com/telegram-mini-apps/init-data-golang"
	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/agurinov/gopl/telegram"
	pl_testing "github.com/agurinov/gopl/testing"
)

var (
	dummyInitData = initdata.User{
		ID:              100500,
		Username:        "johndoe",
		FirstName:       "John",
		LastName:        "Doe",
		IsBot:           false,
		AllowsWriteToPm: true,
	}
	dummyUser = telegram.User{
		ID:           100500,
		Username:     "johndoe",
		FirstName:    "John",
		LastName:     "Doe",
		IsBot:        false,
		AuthorityBot: "DummyBot",
		PrivateChat: telegram.PrivateChat{
			Enabled: true,
			ID:      100500,
		},
	}
)

func hashTmaToken(
	t *testing.T,
	botToken string,
) string {
	t.Helper()

	user := dummyInitData

	var b bytes.Buffer

	require.NoError(t,
		json.NewEncoder(&b).Encode(&user),
	)

	q := url.Values{}
	q.Set("user", b.String())

	hash, err := initdata.SignQueryString(q.Encode(), botToken, time.Now())
	require.NoError(t, err)

	q.Set("hash", hash)

	return "tma " + q.Encode()
}

func TestAuth_authFunc(t *testing.T) {
	type (
		args struct {
			ctx              context.Context
			request          *http.Request
			noSignatureCheck bool
		}
		results struct {
			httpStatusCode int
			httpContent    string
			grpcStatusCode codes.Code
			grpcOut        any
		}
	)

	fooBotToken := "foobot_token"

	newRequest := func(authHeader string) *http.Request {
		request := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
		request.Header.Set("Authorization", authHeader)

		return request
	}

	httpHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := telegram.GetUser(r.Context())
		if err != nil {
			http.Error(w, "oops", http.StatusInternalServerError)

			return
		}

		io.WriteString(w, user.Username) //nolint:errcheck
	})

	newCtx := func(authHeader string) context.Context {
		md := metadata.MD{}
		md = md.Set("authorization", authHeader)

		return md.ToIncoming(context.Background())
	}

	grpcHandler := func(ctx context.Context, req any) (any, error) {
		user, err := telegram.GetUser(ctx)
		if err != nil {
			return nil, status.Error(codes.Internal, "oops")
		}

		return user, nil
	}

	cases := map[string]struct {
		args    args
		results results
		pl_testing.TestCase
	}{
		"case00: no header": {
			args: args{
				ctx:     newCtx(""),
				request: newRequest(""),
			},
			results: results{
				grpcStatusCode: codes.Unauthenticated,
				grpcOut:        nil,
				httpStatusCode: http.StatusUnauthorized,
				httpContent:    "\n",
			},
		},
		"case01: non auth header format": {
			args: args{
				ctx:     newCtx("foobar"),
				request: newRequest("foobar"),
			},
			results: results{
				grpcStatusCode: codes.Unauthenticated,
				grpcOut:        nil,
				httpStatusCode: http.StatusUnauthorized,
				httpContent:    "\n",
			},
		},
		"case02: wrong schema": {
			args: args{
				ctx:     newCtx("bearer foobar"),
				request: newRequest("bearer foobar"),
			},
			results: results{
				grpcStatusCode: codes.Unauthenticated,
				grpcOut:        nil,
				httpStatusCode: http.StatusUnauthorized,
				httpContent:    "\n",
			},
		},
		"case03: wrong token": {
			args: args{
				ctx:     newCtx("tma foobar"),
				request: newRequest("tma foobar"),
			},
			results: results{
				grpcStatusCode: codes.Unauthenticated,
				grpcOut:        nil,
				httpStatusCode: http.StatusUnauthorized,
				httpContent:    "\n",
			},
		},
		"case04: wrong hash": {
			args: args{
				ctx:     newCtx(hashTmaToken(t, "invalid_bot_token")),
				request: newRequest(hashTmaToken(t, "invalid_bot_token")),
			},
			results: results{
				grpcStatusCode: codes.Unauthenticated,
				grpcOut:        nil,
				httpStatusCode: http.StatusUnauthorized,
				httpContent:    "\n",
			},
		},
		"case05: wright hash": {
			args: args{
				ctx:     newCtx(hashTmaToken(t, fooBotToken)),
				request: newRequest(hashTmaToken(t, fooBotToken)),
			},
			results: results{
				grpcStatusCode: codes.OK,
				grpcOut:        dummyUser,
				httpStatusCode: http.StatusOK,
				httpContent:    "johndoe",
			},
			TestCase: pl_testing.TestCase{
				Skip: true,
			},
		},
		"case06: dummy: no header": {
			args: args{
				ctx:              newCtx(""),
				request:          newRequest(""),
				noSignatureCheck: true,
			},
			results: results{
				grpcStatusCode: codes.Unauthenticated,
				grpcOut:        nil,
				httpStatusCode: http.StatusUnauthorized,
				httpContent:    "\n",
			},
		},
		"case07: dummy: non auth header format": {
			args: args{
				ctx:              newCtx("foobar"),
				request:          newRequest("foobar"),
				noSignatureCheck: true,
			},
			results: results{
				grpcStatusCode: codes.Unauthenticated,
				grpcOut:        nil,
				httpStatusCode: http.StatusUnauthorized,
				httpContent:    "\n",
			},
		},
		"case08: dummy: wrong schema": {
			args: args{
				ctx:              newCtx("bearer foobar"),
				request:          newRequest("bearer foobar"),
				noSignatureCheck: true,
			},
			results: results{
				grpcStatusCode: codes.Unauthenticated,
				grpcOut:        nil,
				httpStatusCode: http.StatusUnauthorized,
				httpContent:    "\n",
			},
		},
		"case09: dummy: wrong token": {
			args: args{
				ctx:              newCtx("tma foobar"),
				request:          newRequest("tma foobar"),
				noSignatureCheck: true,
			},
			results: results{
				grpcStatusCode: codes.Unauthenticated,
				grpcOut:        nil,
				httpStatusCode: http.StatusUnauthorized,
				httpContent:    "\n",
			},
		},
		"case10: dummy: wrong hash": {
			args: args{
				ctx:              newCtx(hashTmaToken(t, "invalid_bot_token")),
				request:          newRequest(hashTmaToken(t, "invalid_bot_token")),
				noSignatureCheck: true,
			},
			results: results{
				grpcStatusCode: codes.OK,
				grpcOut:        dummyUser,
				httpStatusCode: http.StatusOK,
				httpContent:    "johndoe",
			},
		},
		"case11: dummy: wright hash": {
			args: args{
				ctx:              newCtx(hashTmaToken(t, fooBotToken)),
				request:          newRequest(hashTmaToken(t, fooBotToken)),
				noSignatureCheck: true,
			},
			results: results{
				grpcStatusCode: codes.OK,
				grpcOut:        dummyUser,
				httpStatusCode: http.StatusOK,
				httpContent:    "johndoe",
			},
		},
	}

	for name := range cases {
		tc := cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			auth, err := telegram.NewAuth(
				telegram.WithAuthLogger(zaptest.NewLogger(t)),
				telegram.WithAuthNoSignatureCheck(tc.args.noSignatureCheck),
				telegram.WithAuthBotTokens(map[string]string{"FooBot": fooBotToken}),
			)
			require.NoError(t, err)
			require.NotNil(t, auth)

			grpcOut, err := auth.UnaryServerInterceptor(
				tc.args.ctx,
				nil,
				new(grpc.UnaryServerInfo),
				grpcHandler,
			)
			grpcCode := status.Code(err)

			require.Equal(t, tc.results.grpcStatusCode, grpcCode)
			require.Equal(t, tc.results.grpcOut, grpcOut)

			var (
				recorder    = httptest.NewRecorder()
				authHandler = auth.Middleware(httpHandler)
			)

			authHandler.ServeHTTP(recorder, tc.args.request)

			require.Equal(t, tc.results.httpStatusCode, recorder.Code)
			require.Equal(t, tc.results.httpContent, recorder.Body.String())
		})
	}
}

func TestAuth_Validate(t *testing.T) {
	type (
		args struct {
			botTokens map[string]string
		}
	)

	cases := map[string]struct {
		args args
		pl_testing.TestCase
	}{
		"case00: no bot tokens": {
			args: args{},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailAsErr: new(validator.ValidationErrors),
			},
		},
		"case01: empty key": {
			args: args{
				botTokens: map[string]string{
					"": "token",
				},
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailAsErr: new(validator.ValidationErrors),
			},
		},
		"case02: empty token": {
			args: args{
				botTokens: map[string]string{
					"FooBot": "",
				},
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailAsErr: new(validator.ValidationErrors),
			},
		},
		"case03: success": {
			args: args{
				botTokens: map[string]string{
					"FooBot": "token",
				},
			},
		},
	}

	for name := range cases {
		tc := cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			auth, err := telegram.NewAuth(
				telegram.WithAuthLogger(zaptest.NewLogger(t)),
				telegram.WithAuthBotTokens(tc.args.botTokens),
			)
			tc.CheckError(t, err)
			require.NotNil(t, auth)

			err = auth.Validate()
			tc.CheckError(t, err)
		})
	}
}
