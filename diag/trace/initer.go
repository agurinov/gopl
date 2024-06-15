package trace

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"

	c "github.com/agurinov/gopl/patterns/creational"
)

type (
	initer struct {
		cmdName         string
		exporterOptions []otlptracegrpc.Option
		samplerOptions  []SamplerOption
		batcherOptions  []trace.BatchSpanProcessorOption
	}
	IniterOption = c.Option[initer]
)

func Init(ctx context.Context, opts ...IniterOption) error {
	i, err := c.NewWithValidate(opts...)
	if err != nil {
		return err
	}

	exporter, err := otlptracegrpc.New(ctx, i.exporterOptions...)
	if err != nil {
		return err
	}

	sampler, err := NewSampler(i.samplerOptions...)
	if err != nil {
		return err
	}

	provider := trace.NewTracerProvider(
		trace.WithBatcher(exporter, i.batcherOptions...),
		trace.WithSampler(sampler),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceName(i.cmdName),
			),
		),
	)

	otel.SetTextMapPropagator(propagation.TraceContext{})
	otel.SetTracerProvider(provider)

	return nil
}
