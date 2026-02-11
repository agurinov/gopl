package kafka

import (
	"context"
	"sync"

	"github.com/agurinov/gopl/x"
)

type (
	partitionContext struct {
		ctx        context.Context
		cancelFunc context.CancelFunc
	}
	partitionHolder struct {
		assigned   map[int32]partitionContext
		assignedMu sync.RWMutex
	}
)

// NOTE: idempotent method.
func (h *partitionHolder) revokePartition(
	partition int32,
) {
	h.assignedMu.Lock()
	defer h.assignedMu.Unlock()

	if cancelFunc := h.assigned[partition].cancelFunc; cancelFunc != nil {
		cancelFunc()
	}

	delete(h.assigned, partition)
}

// NOTE: not idempotent method (cancelling previous context).
// SHOULD be used with goroutine spawner.
func (h *partitionHolder) assignPartitions(
	ctx context.Context,
	partitions []int32,
) {
	h.assignedMu.Lock()
	defer h.assignedMu.Unlock()

	for _, p := range partitions {
		if previousCancelFunc := h.assigned[p].cancelFunc; previousCancelFunc != nil {
			previousCancelFunc()
		}

		pCtx, cancelFunc := context.WithCancel(ctx)

		h.assigned[p] = partitionContext{
			ctx:        pCtx,
			cancelFunc: cancelFunc,
		}
	}
}

func (h *partitionHolder) assignedPartitions() []int32 {
	h.assignedMu.RLock()
	defer h.assignedMu.RUnlock()

	return x.MapKeys(h.assigned)
}

func (h *partitionHolder) partitionContext(
	partition int32,
) context.Context {
	h.assignedMu.RLock()
	defer h.assignedMu.RUnlock()

	return h.assigned[partition].ctx
}
