package trace

/*
import (
	"time"

	"go.opentelemetry.io/otel/sdk/trace"
)

// func Providerglobal ???
func NewTraceProvider(
	cmdName string,
	traceRatio float64,
	exporter trace.SpanExporter,
) *trace.TracerProvider {
	return trace.NewTracerProvider(
		trace.WithBatcher(exporter, trace.WithBatchTimeout(5*time.Second)),
		trace.WithSampler(trace.ParentBased(trace.TraceIDRatioBased(traceRatio))),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(cmdName),
			),
		),
	)
}
*/
