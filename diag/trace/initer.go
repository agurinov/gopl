package trace

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"

	"github.com/agurinov/gopl/env/envvars"
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
	if !envvars.OtelTraceEnabled.Present() {
		return nil
	}

	switch enabled, err := envvars.OtelTraceEnabled.Value(); {
	case err != nil:
		return err
	case !enabled:
		return nil
	}

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

	spanProcessor := spanProcessor{
		next:        trace.NewBatchSpanProcessor(exporter),
		ratio:       sampler.sampleRatio,
		minDuration: sampler.sampleDuration,
		errors:      sampler.sampleErrors,
	}

	provider := trace.NewTracerProvider(
		trace.WithBatcher(exporter, i.batcherOptions...),
		trace.WithSpanProcessor(spanProcessor),
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
