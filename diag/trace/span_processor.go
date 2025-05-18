package trace

import (
	"context"
	"encoding/binary"
	"time"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/sdk/trace"
)

type (
	spanProcessor struct {
		next        trace.SpanProcessor
		minDuration time.Duration
		errors      bool
		ratio       float64
	}
)

func (sp spanProcessor) OnStart(parent context.Context, s trace.ReadWriteSpan) {
	sp.next.OnStart(parent, s)
}

func (sp spanProcessor) Shutdown(ctx context.Context) error {
	return sp.next.Shutdown(ctx)
}

func (sp spanProcessor) ForceFlush(ctx context.Context) error {
	return sp.next.ForceFlush(ctx)
}

func (sp spanProcessor) OnEnd(s trace.ReadOnlySpan) {
	if sp.errors && s.Status().Code == codes.Error {
		sp.next.OnEnd(s)
	}

	var (
		spanDuration = s.EndTime().Sub(s.StartTime())
		isLongSpan   = spanDuration >= sp.minDuration
	)

	if sp.minDuration != 0 && isLongSpan {
		sp.next.OnEnd(s)
	}

	var (
		ctx               = s.SpanContext()
		traceID           = ctx.TraceID()
		x                 = binary.BigEndian.Uint64(traceID[8:16]) >> 1
		traceIDUpperBound = uint64(sp.ratio * (1 << 63)) //nolint:gomnd,mnd
	)

	if x < traceIDUpperBound {
		sp.next.OnEnd(s)
	}
}
