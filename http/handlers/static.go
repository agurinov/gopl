package handlers

import (
	"io/fs"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"

	"github.com/agurinov/gopl/diag/metrics"
	"github.com/agurinov/gopl/http/middlewares"
	c "github.com/agurinov/gopl/patterns/creational"
)

type (
	static struct {
		fs                fs.FS
		logger            *zap.Logger
		knownPaths        map[string][]byte
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
		middlewares.AccessLog(h.logger),
		chimw.GetHead,
	)

	r.Use(h.customMiddlewares...)
	// r.Use(middlewares.Panic(obj.logger))

	fsHandler := http.FileServer(http.FS(h.fs))

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

		switch f, isKnownPath := h.knownPaths[path]; isKnownPath {
		case true:
			mimeType := mime.TypeByExtension(ext)
			if mimeType == "" {
				mimeType = http.DetectContentType(f)
			}

			w.Header().Set("Content-Type", mimeType)
			w.Write(f) //nolint:errcheck
		default:
			fsHandler.ServeHTTP(w, r)
		}
	})

	return r
}
