package sql

import (
	"fmt"

	"go.opentelemetry.io/otel/trace"
)

type Query string

const sqlQueryWrapper = `%s
/*
traceparent='%s'
*/`

func (q Query) String() string { return string(q) }

func (q Query) WithSpan(span trace.Span) string {
	var (
		ctx     = span.SpanContext()
		traceid = ctx.TraceID().String()
	)

	return fmt.Sprintf(
		sqlQueryWrapper,
		q.String(),
		traceid,
	)
}
