package handlers

import (
	"go.uber.org/zap"

	"github.com/agurinov/gopl/diag/probes"
)

func WithDebugAtomicLevel(lvl *zap.AtomicLevel) DebugOption {
	return func(h *debug) error {
		h.atomicLogLevel = lvl

		return nil
	}
}

func WithDebugLogger(logger *zap.Logger) DebugOption {
	return func(h *debug) error {
		if logger == nil {
			return nil
		}

		h.logger = logger.Named("http.handler.debug")

		return nil
	}
}

func WithDebugProber(prober *probes.Prober) DebugOption {
	return func(h *debug) error {
		h.prober = prober

		return nil
	}
}
