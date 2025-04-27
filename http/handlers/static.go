package handlers

import (
	"context"
	"io/fs"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/exp/slices"

	"github.com/agurinov/gopl/diag/metrics"
	"github.com/agurinov/gopl/http/middlewares"
	c "github.com/agurinov/gopl/patterns/creational"
)

type (
	bufferFunc func(context.Context) ([]byte, error)
	static     struct {
		fs                fs.FS
		logger            *zap.Logger
		knownBufFunc      map[string]bufferFunc
		noCachePaths      []string
		customMiddlewares []middlewares.Middleware
		spaEnabled        bool
	}
	StaticOption c.Option[static]
)

var NewStatic = c.NewWithValidate[static, StaticOption]

func (h static) Handler() http.Handler {
	r := chi.NewRouter()

	r.Use(
		middlewares.Trace,
		middlewares.Metrics(
			metrics.WithBuckets(metrics.BucketFast),
		),
		middlewares.AccessLog(
			h.logger,
			zapcore.InfoLevel,
		),
		chimw.GetHead,
	)

	r.Use(h.customMiddlewares...)

	bundle := http.FileServer(http.FS(h.fs))

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		var (
			path  = filepath.Clean(r.URL.Path)
			ext   = filepath.Ext(path)
			isDir = ext == ""
		)

		if h.spaEnabled && isDir {
			path = "/"
		}

		if slices.Contains(h.noCachePaths, path) {
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")
		}

		r.URL.Path = path

		bufGetter, isKnownPath := h.knownBufFunc[path]

		switch {
		case isKnownPath:
			var (
				ctx      = r.Context()
				buf, err = bufGetter(ctx)
			)

			if err != nil {
				h.logger.Error(
					"can't get buffer for known path",
					zap.String("path", path),
					zap.Error(err),
				)

				w.WriteHeader(http.StatusInternalServerError)

				break
			}

			mimeType := mime.TypeByExtension(ext)
			if mimeType == "" {
				mimeType = http.DetectContentType(buf)
			}

			w.Header().Set("Content-Type", mimeType)
			w.Write(buf) //nolint:errcheck
		default:
			bundle.ServeHTTP(w, r)
		}
	})

	return r
}
