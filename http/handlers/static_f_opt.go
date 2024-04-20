package handlers

import (
	"io/fs"

	"go.uber.org/zap"
)

func WithStaticLogger(logger *zap.Logger) StaticOption {
	return func(h *static) error {
		h.logger = logger.Named("http.handler.static")

		return nil
	}
}

func WithStaticFS(staticFS fs.FS, dirname string) StaticOption {
	return func(h *static) error {
		rootFS, err := fs.Sub(staticFS, dirname)
		if err != nil {
			return err
		}

		h.fs = rootFS

		return nil
	}
}

func WithStaticSPA(spaEnabled bool) StaticOption {
	return func(h *static) error {
		h.spaEnabled = spaEnabled

		return nil
	}
}
