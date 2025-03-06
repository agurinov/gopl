package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/agurinov/gopl/http/handlers"
	"github.com/agurinov/gopl/nopanic"
	pl_testing "github.com/agurinov/gopl/testing"
)

func TestBasic(t *testing.T) {
	type (
		args struct {
			request             *http.Request
			basicHandlerOptions []handlers.BasicOption
		}
		results struct {
			headers    http.Header
			statusCode int
		}
	)

	nopanicHandler, err := nopanic.NewHandler(
		nopanic.WithLogger(zaptest.NewLogger(t)),
	)
	require.NoError(t, err)
	require.NotNil(t, nopanicHandler)

	cases := map[string]struct {
		pl_testing.TestCase
		results results
		args    args
	}{
		"case00: check get": {
			args: args{
				basicHandlerOptions: []handlers.BasicOption{
					handlers.WithBasicHandlers(map[string]http.Handler{
						"/*": http.HandlerFunc(
							func(w http.ResponseWriter, r *http.Request) {
								http.Redirect(w, r, "/foobar", http.StatusPermanentRedirect)
							},
						),
					}),
				},
				request: httptest.NewRequest(http.MethodGet, "/", nil),
			},
			results: results{
				statusCode: http.StatusPermanentRedirect,
				headers: http.Header{
					"Location":     []string{"/foobar"},
					"Content-Type": []string{"text/html; charset=utf-8"},
				},
			},
		},
		"case01: check post": {
			args: args{
				basicHandlerOptions: []handlers.BasicOption{
					handlers.WithBasicHandlers(map[string]http.Handler{
						"/*": http.HandlerFunc(
							func(w http.ResponseWriter, r *http.Request) {
								http.Redirect(w, r, "/foobar", http.StatusPermanentRedirect)
							},
						),
					}),
				},
				request: httptest.NewRequest(http.MethodPost, "/", nil),
			},
			results: results{
				statusCode: http.StatusPermanentRedirect,
				headers: http.Header{
					"Location": []string{"/foobar"},
				},
			},
		},
		"case03: check put foobar": {
			args: args{
				basicHandlerOptions: []handlers.BasicOption{
					handlers.WithBasicHandlers(map[string]http.Handler{
						"/*": http.HandlerFunc(
							func(w http.ResponseWriter, r *http.Request) {
								http.Redirect(w, r, "/foobar/baz", http.StatusPermanentRedirect)
							},
						),
					}),
				},
				request: httptest.NewRequest(http.MethodPut, "/foobar", nil),
			},
			results: results{
				statusCode: http.StatusPermanentRedirect,
				headers: http.Header{
					"Location": []string{"/foobar/baz"},
				},
			},
		},
		"case04: panic in handler": {
			args: args{
				basicHandlerOptions: []handlers.BasicOption{
					handlers.WithBasicHandlers(map[string]http.Handler{
						"/*": http.HandlerFunc(
							func(_ http.ResponseWriter, _ *http.Request) {
								panic("OOPS")
							},
						),
					}),
				},
				request: httptest.NewRequest(http.MethodPut, "/foobar", nil),
			},
			results: results{
				statusCode: http.StatusInternalServerError,
				headers:    http.Header{},
			},
		},
	}

	for name := range cases {
		name, tc := name, cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			recorder := httptest.NewRecorder()

			tc.args.basicHandlerOptions = append(tc.args.basicHandlerOptions,
				handlers.WithBasicLogger(zaptest.NewLogger(t)),
				handlers.WithBasicCustomMiddlewares(
					nopanicHandler.Middleware(),
				),
			)

			basic, err := handlers.NewBasic(tc.args.basicHandlerOptions...)
			require.NoError(t, err)
			require.NotNil(t, basic)

			handler := basic.Handler()
			require.NotNil(t, handler)

			handler.ServeHTTP(recorder, tc.args.request)

			require.Equal(t, tc.results.statusCode, recorder.Code)
			require.Equal(t, tc.results.headers, recorder.Header())

			require.Nil(t, recover())
		})
	}
}
