package handlers

import "go.uber.org/zap"

func WithAtomicLevel(lvl *zap.AtomicLevel) DebugOption {
	return func(h *debug) error {
		h.atomicLogLevel = lvl

		return nil
	}
}
