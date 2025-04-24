package handlers

import (
	"context"
	"io/fs"
	"path/filepath"

	"go.uber.org/zap"

	"github.com/agurinov/gopl/http/middlewares"
)

func WithStaticLogger(logger *zap.Logger) StaticOption {
	return func(h *static) error {
		if logger == nil {
			return nil
		}

		h.logger = logger.Named("http.handler.static")

		return nil
	}
}

func WithStaticBundle(staticFS fs.FS, dirname string) StaticOption {
	return func(h *static) error {
		if dirname != "" {
			subFS, err := fs.Sub(staticFS, dirname)
			if err != nil {
				return err
			}

			staticFS = subFS
		}

		h.fs = staticFS

		return nil
	}
}

func WithStaticKnownBuffer(path string, buf []byte) StaticOption {
	return WithStaticKnownBufferFunc(
		path,
		func(context.Context) ([]byte, error) {
			return buf, nil
		},
	)
}

func WithStaticKnownBufferFunc(path string, f bufferFunc) StaticOption {
	return func(h *static) error {
		path = filepath.Clean(path)

		if h.knownBufFunc == nil {
			h.knownBufFunc = make(map[string]bufferFunc)
		}

		h.knownBufFunc[path] = f

		return nil
	}
}

func WithStaticSPA(spaEnabled bool) StaticOption {
	return func(h *static) error {
		h.spaEnabled = spaEnabled

		return nil
	}
}

func WithStaticNoCachePaths(paths ...string) StaticOption {
	return func(h *static) error {
		h.noCachePaths = paths

		return nil
	}
}

func WithStaticCustomMiddlewares(mw ...middlewares.Middleware) StaticOption {
	return func(s *static) error {
		s.customMiddlewares = append(s.customMiddlewares, mw...)

		return nil
	}
}
