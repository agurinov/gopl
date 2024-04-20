package handlers_test

import (
	"embed"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/agurinov/gopl/http/handlers"
	pl_testing "github.com/agurinov/gopl/testing"
)

//go:embed all:testdata
var content embed.FS

func TestStatic(t *testing.T) {
	type (
		args struct {
			staticOptions []handlers.StaticOption
			request       *http.Request
		}
		results struct {
			statusCode int
			content    string
		}
	)

	cases := map[string]struct {
		args    args
		results results
		pl_testing.TestCase
	}{
		"case00: index on slash": {
			args: args{
				staticOptions: []handlers.StaticOption{
					handlers.WithStaticLogger(zaptest.NewLogger(t)),
					handlers.WithStaticFS(content, "testdata"),
				},
				request: httptest.NewRequest(http.MethodGet, "/", nil),
			},
			results: results{
				statusCode: http.StatusOK,
				content:    "index.html\n",
			},
		},
		"case01: static file 200": {
			args: args{
				staticOptions: []handlers.StaticOption{
					handlers.WithStaticLogger(zaptest.NewLogger(t)),
					handlers.WithStaticFS(content, "testdata"),
				},
				request: httptest.NewRequest(http.MethodGet, "/robots.txt", nil),
			},
			results: results{
				statusCode: http.StatusOK,
				content:    "robots.txt\n",
			},
		},
		"case02: static file 404": {
			args: args{
				staticOptions: []handlers.StaticOption{
					handlers.WithStaticLogger(zaptest.NewLogger(t)),
					handlers.WithStaticFS(content, "testdata"),
				},
				request: httptest.NewRequest(http.MethodGet, "/js/foo.txt", nil),
			},
			results: results{
				statusCode: http.StatusNotFound,
				content:    "404 page not found\n",
			},
		},
		"case03: non GET not allowed": {
			args: args{
				staticOptions: []handlers.StaticOption{
					handlers.WithStaticLogger(zaptest.NewLogger(t)),
					handlers.WithStaticFS(content, "testdata"),
				},
				request: httptest.NewRequest(http.MethodPost, "/js/foo.txt", nil),
			},
			results: results{
				statusCode: http.StatusMethodNotAllowed,
			},
		},
		"case04: SPA static 200": {
			args: args{
				staticOptions: []handlers.StaticOption{
					handlers.WithStaticLogger(zaptest.NewLogger(t)),
					handlers.WithStaticFS(content, "testdata"),
					handlers.WithStaticSPA(true),
				},
				request: httptest.NewRequest(http.MethodGet, "/js/main.js", nil),
			},
			results: results{
				statusCode: http.StatusOK,
				content:    "main.js\n",
			},
		},
		"case05: SPA static 404": {
			args: args{
				staticOptions: []handlers.StaticOption{
					handlers.WithStaticLogger(zaptest.NewLogger(t)),
					handlers.WithStaticFS(content, "testdata"),
					handlers.WithStaticSPA(true),
				},
				request: httptest.NewRequest(http.MethodGet, "/js/foo.js/", nil),
			},
			results: results{
				statusCode: http.StatusNotFound,
				content:    "404 page not found\n",
			},
		},
		"case06: SPA any other route on slash": {
			args: args{
				staticOptions: []handlers.StaticOption{
					handlers.WithStaticLogger(zaptest.NewLogger(t)),
					handlers.WithStaticFS(content, "testdata"),
					handlers.WithStaticSPA(true),
				},
				request: httptest.NewRequest(http.MethodGet, "/js", nil),
			},
			results: results{
				statusCode: http.StatusOK,
				content:    "index.html\n",
			},
		},
	}

	for name := range cases {
		name, tc := name, cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			recorder := httptest.NewRecorder()

			static, err := handlers.NewStatic(tc.args.staticOptions...)
			require.NoError(t, err)
			require.NotNil(t, static)

			handler := static.Handler()
			require.NotNil(t, handler)

			handler.ServeHTTP(recorder, tc.args.request)

			require.Equal(t, tc.results.statusCode, recorder.Code)
			require.Equal(t, tc.results.content, recorder.Body.String())
		})
	}
}
