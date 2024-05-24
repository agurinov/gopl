package trace

import (
	"context"

	"go.opentelemetry.io/otel/propagation"
)

func TraceparentFromContext(ctx context.Context) string {
	tc := propagation.TraceContext{}
	mc := propagation.MapCarrier{}

	tc.Inject(ctx, mc)

	for _, k := range mc.Keys() {
		if k == "traceparent" {
			return mc.Get(k)
		}
	}

	return ""
}

func TraceparentToContext(ctx context.Context, traceparent string) context.Context {
	if traceparent == "" {
		return ctx
	}

	tc := propagation.TraceContext{}
	mc := propagation.MapCarrier{}

	mc.Set("traceparent", traceparent)

	return tc.Extract(ctx, mc)
}
