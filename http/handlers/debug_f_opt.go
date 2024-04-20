package handlers

import "go.uber.org/zap"

func WithDebugAtomicLevel(lvl *zap.AtomicLevel) DebugOption {
	return func(h *debug) error {
		h.atomicLogLevel = lvl

		return nil
	}
}

func WithDebugLogger(logger *zap.Logger) DebugOption {
	return func(h *debug) error {
		h.logger = logger.Named("http.handler.debug")

		return nil
	}
}
