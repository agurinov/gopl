package trace

import (
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/trace"
)

func WithCmdName(cmdName string) IniterOption {
	return func(i *initer) error {
		i.cmdName = cmdName

		return nil
	}
}

func WithBatcherOptions(opts ...trace.BatchSpanProcessorOption) IniterOption {
	return func(i *initer) error {
		i.batcherOptions = opts

		return nil
	}
}

func WithSamplerOptions(opts ...SamplerOption) IniterOption {
	return func(i *initer) error {
		i.samplerOptions = opts

		return nil
	}
}

func WithExporterOptions(opts ...otlptracegrpc.Option) IniterOption {
	return func(i *initer) error {
		i.exporterOptions = opts

		return nil
	}
}
