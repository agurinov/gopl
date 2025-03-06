package handlers_test

import (
	"embed"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/agurinov/gopl/http/handlers"
	"github.com/agurinov/gopl/nopanic"
	pl_testing "github.com/agurinov/gopl/testing"
)

//go:embed testdata
var bundle embed.FS

func TestStatic(t *testing.T) {
	type (
		args struct {
			request              *http.Request
			staticHandlerOptions []handlers.StaticOption
		}
		results struct {
			headers    http.Header
			content    string
			statusCode int
		}
	)

	nopanicHandler, err := nopanic.NewHandler(
		nopanic.WithLogger(zaptest.NewLogger(t)),
	)
	require.NoError(t, err)
	require.NotNil(t, nopanicHandler)

	cases := map[string]struct {
		results results
		pl_testing.TestCase
		args args
	}{
		"case00: embed: index on slash": {
			args: args{
				staticHandlerOptions: []handlers.StaticOption{
					handlers.WithStaticBundle(bundle, "testdata"),
				},
				request: httptest.NewRequest(http.MethodGet, "/", nil),
			},
			results: results{
				statusCode: http.StatusOK,
				content:    "index.html\n",
				headers: http.Header{
					"Accept-Ranges":  []string{"bytes"},
					"Content-Length": []string{"11"},
					"Content-Type":   []string{"text/html; charset=utf-8"},
				},
			},
		},
		"case01: embed: static file 200": {
			args: args{
				staticHandlerOptions: []handlers.StaticOption{
					handlers.WithStaticBundle(bundle, "testdata"),
				},
				request: httptest.NewRequest(http.MethodGet, "/robots.txt", nil),
			},
			results: results{
				statusCode: http.StatusOK,
				content:    "robots.txt\n",
				headers: http.Header{
					"Accept-Ranges":  []string{"bytes"},
					"Content-Length": []string{"11"},
					"Content-Type":   []string{"text/plain; charset=utf-8"},
				},
			},
		},
		"case02: embed: static file 404": {
			args: args{
				staticHandlerOptions: []handlers.StaticOption{
					handlers.WithStaticBundle(bundle, "testdata"),
				},
				request: httptest.NewRequest(http.MethodGet, "/js/foo.txt", nil),
			},
			results: results{
				statusCode: http.StatusNotFound,
				content:    "404 page not found\n",
				headers: http.Header{
					"Accept-Ranges":  []string{"bytes"},
					"Content-Length": []string{"11"},
					"Content-Type":   []string{"text/html; charset=utf-8"},
				},
			},
		},
		"case03: embed: nothing except GET is allowed": {
			args: args{
				staticHandlerOptions: []handlers.StaticOption{
					handlers.WithStaticBundle(bundle, "testdata"),
				},
				request: httptest.NewRequest(http.MethodPost, "/js/foo.txt", nil),
			},
			results: results{
				statusCode: http.StatusMethodNotAllowed,
				headers: http.Header{
					"Allow": []string{"GET"},
				},
			},
		},
		"case04: embed: SPA: static file 200": {
			args: args{
				staticHandlerOptions: []handlers.StaticOption{
					handlers.WithStaticBundle(bundle, "testdata"),
					handlers.WithStaticSPA(true),
				},
				request: httptest.NewRequest(http.MethodGet, "/js/main.js", nil),
			},
			results: results{
				statusCode: http.StatusOK,
				content:    "main.js\n",
				headers: http.Header{
					"Accept-Ranges":  []string{"bytes"},
					"Content-Length": []string{"8"},
					"Content-Type":   []string{"text/javascript; charset=utf-8"},
				},
			},
		},
		"case05: embed: SPA: static file 404": {
			args: args{
				staticHandlerOptions: []handlers.StaticOption{
					handlers.WithStaticBundle(bundle, "testdata"),
					handlers.WithStaticSPA(true),
				},
				request: httptest.NewRequest(http.MethodGet, "/js/foo.js/", nil),
			},
			results: results{
				statusCode: http.StatusNotFound,
				content:    "404 page not found\n",
			},
		},
		"case06: embed: SPA: any route on slash": {
			args: args{
				staticHandlerOptions: []handlers.StaticOption{
					handlers.WithStaticBundle(bundle, "testdata"),
					handlers.WithStaticSPA(true),
				},
				request: httptest.NewRequest(http.MethodGet, "/js", nil),
			},
			results: results{
				statusCode: http.StatusOK,
				content:    "index.html\n",
				headers: http.Header{
					"Accept-Ranges":  []string{"bytes"},
					"Content-Length": []string{"11"},
					"Content-Type":   []string{"text/html; charset=utf-8"},
				},
			},
		},
		"case07: os fs: SPA any route on slash": {
			args: args{
				staticHandlerOptions: []handlers.StaticOption{
					handlers.WithStaticBundle(os.DirFS("testdata"), ""),
					handlers.WithStaticSPA(true),
				},
				request: httptest.NewRequest(http.MethodGet, "/js", nil),
			},
			results: results{
				statusCode: http.StatusOK,
				content:    "index.html\n",
				headers: http.Header{
					"Accept-Ranges":  []string{"bytes"},
					"Content-Length": []string{"11"},
					"Content-Type":   []string{"text/html; charset=utf-8"},
				},
			},
		},
		"case08: os fs: SPA static file 404": {
			args: args{
				staticHandlerOptions: []handlers.StaticOption{
					handlers.WithStaticBundle(os.DirFS("testdata"), ""),
					handlers.WithStaticSPA(true),
				},
				request: httptest.NewRequest(http.MethodGet, "/js/foo.js/", nil),
			},
			results: results{
				statusCode: http.StatusNotFound,
				content:    "404 page not found\n",
			},
		},
		"case09: embed: static file 200 non cacheable": {
			args: args{
				staticHandlerOptions: []handlers.StaticOption{
					handlers.WithStaticBundle(bundle, "testdata"),
					handlers.WithStaticNoCachePaths("/", "/index.html", "/robots.txt"),
				},
				request: httptest.NewRequest(http.MethodGet, "/robots.txt", nil),
			},
			results: results{
				statusCode: http.StatusOK,
				content:    "robots.txt\n",
				headers: http.Header{
					"Accept-Ranges":  []string{"bytes"},
					"Content-Length": []string{"11"},
					"Content-Type":   []string{"text/plain; charset=utf-8"},
					"Cache-Control":  []string{"no-cache, no-store, must-revalidate"},
					"Pragma":         []string{"no-cache"},
					"Expires":        []string{"0"},
				},
			},
		},
		"case10: embed: SPA any route on slash non cacheable": {
			args: args{
				staticHandlerOptions: []handlers.StaticOption{
					handlers.WithStaticBundle(bundle, "testdata"),
					handlers.WithStaticSPA(true),
					handlers.WithStaticNoCachePaths("/", "/index.html"),
				},
				request: httptest.NewRequest(http.MethodGet, "/foo/bar", nil),
			},
			results: results{
				statusCode: http.StatusOK,
				content:    "index.html\n",
				headers: http.Header{
					"Accept-Ranges":  []string{"bytes"},
					"Content-Length": []string{"11"},
					"Content-Type":   []string{"text/html; charset=utf-8"},
					"Cache-Control":  []string{"no-cache, no-store, must-revalidate"},
					"Pragma":         []string{"no-cache"},
					"Expires":        []string{"0"},
				},
			},
		},
		"case11: embed: SPA: known file 200": {
			args: args{
				staticHandlerOptions: []handlers.StaticOption{
					handlers.WithStaticBundle(bundle, "testdata"),
					handlers.WithStaticKnownFile("/config.json", []byte(`{"foo":"bar"}`)),
					handlers.WithStaticSPA(true),
					handlers.WithStaticNoCachePaths("/config.json"),
				},
				request: httptest.NewRequest(http.MethodGet, "/config.json", nil),
			},
			results: results{
				statusCode: http.StatusOK,
				content:    `{"foo":"bar"}`,
				headers: http.Header{
					"Content-Type":  []string{"application/json"},
					"Cache-Control": []string{"no-cache, no-store, must-revalidate"},
					"Pragma":        []string{"no-cache"},
					"Expires":       []string{"0"},
				},
			},
		},
		"case12: middleware: use custom mw": {
			args: args{
				staticHandlerOptions: []handlers.StaticOption{
					handlers.WithStaticBundle(bundle, "testdata"),
					handlers.WithStaticKnownFile("/config.json", []byte(`{"foo":"bar"}`)),
					handlers.WithStaticSPA(true),
					handlers.WithStaticNoCachePaths("/config.json"),
					handlers.WithStaticCustomMiddlewares(func(h http.Handler) http.Handler {
						return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							w.Header().Set("Custom-Header", "Value")
						})
					}),
				},
				request: httptest.NewRequest(http.MethodGet, "/", nil),
			},
			results: results{
				statusCode: http.StatusOK,
				headers: http.Header{
					"Custom-Header": []string{"Value"},
				},
			},
		},
	}

	for name := range cases {
		name, tc := name, cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			recorder := httptest.NewRecorder()

			tc.args.staticHandlerOptions = append(tc.args.staticHandlerOptions,
				handlers.WithStaticLogger(zaptest.NewLogger(t)),
				handlers.WithStaticCustomMiddlewares(
					nopanicHandler.Middleware(),
				),
			)

			static, err := handlers.NewStatic(tc.args.staticHandlerOptions...)
			require.NoError(t, err)
			require.NotNil(t, static)

			handler := static.Handler()
			require.NotNil(t, handler)

			handler.ServeHTTP(recorder, tc.args.request)

			require.Equal(t, tc.results.statusCode, recorder.Code)
			require.Equal(t, tc.results.content, recorder.Body.String())

			if tc.results.statusCode == http.StatusNotFound {
				return
			}

			headers := recorder.Header()

			delete(headers, "Last-Modified")

			require.Equal(t, tc.results.headers, headers)
		})
	}
}
