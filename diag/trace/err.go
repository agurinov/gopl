package trace

import (
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func RegisterError(span trace.Span, err error) {
	if err != nil {
		span.RecordError(
			err,
			trace.WithStackTrace(true),
		)
		span.SetStatus(codes.Error, err.Error())
	}
}

func CatchError(span trace.Span, err error) error {
	RegisterError(span, err)

	return err
}
