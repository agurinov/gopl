package trace

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func SpanName() string {
	// fnName := GetCallerName(3)
	// return getSpanNameFromCallerName(fnName)
	return "implement me"
}

func StartSpan(
	ctx context.Context,
	spanName string,
	opts ...trace.SpanStartOption,
) (
	context.Context,
	trace.Span,
) {
	var (
		tracerName = ""
		tracer     = trace.SpanFromContext(ctx).TracerProvider().Tracer(tracerName)
	)

	ctx, span := tracer.Start(ctx, spanName, opts...)
	if span.IsRecording() {
		return ctx, span
	}

	return otel.Tracer(tracerName).Start(ctx, spanName, opts...)
}

func StartNamedSpan(
	ctx context.Context,
	opts ...trace.SpanStartOption,
) (
	context.Context,
	trace.Span,
) {
	return StartSpan(ctx, SpanName(), opts...)
}
