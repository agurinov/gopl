package sql

import (
	"fmt"

	"go.opentelemetry.io/otel/trace"
)

type Query string

func (q Query) String() string { return string(q) }

func (q Query) WithSpan(span trace.Span) string {
	var (
		ctx     = span.SpanContext()
		traceid = ctx.TraceID().String()
		spanid  = ctx.SpanID().String()
	)

	return fmt.Sprintf(
		"-- trace_id: %s\n-- span_id: %s\n%s",
		traceid,
		spanid,
		q.String(),
	)
}
