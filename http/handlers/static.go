package handlers

import (
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"

	"github.com/agurinov/gopl/http/middlewares"
	c "github.com/agurinov/gopl/patterns/creational"
)

type (
	static struct {
		logger       *zap.Logger
		fs           fs.FS
		noCachePaths []string
		spaEnabled   bool
	}
	StaticOption c.Option[static]
)

var NewStatic = c.NewWithValidate[static, StaticOption]

func (h static) Handler() http.Handler {
	r := chi.NewRouter()

	r.Use(
		middlewares.Trace,
		// middlewares.Metrics(),
		middlewares.AccessLog(h.logger),
		chimw.GetHead,
		// middlewares.Panic(obj.logger),
	)

	fsHandler := http.FileServer(http.FS(h.fs))

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		if h.spaEnabled {
			var (
				path = filepath.Clean(r.URL.Path)
				ext  = filepath.Ext(path)
			)

			if isDir := ext == ""; isDir {
				r.URL.Path = "/"
			}
		}

		if slices.Contains(h.noCachePaths, r.URL.Path) {
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")
		}

		fsHandler.ServeHTTP(w, r)
	})

	return r
}
