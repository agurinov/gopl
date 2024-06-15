package trace

import (
	"time"

	"go.opentelemetry.io/otel/sdk/trace"
)

func WithSampleError(sampleErrors bool) SamplerOption {
	return func(s *sampler) error {
		s.sampleErrors = sampleErrors

		return nil
	}
}

func WithSampleDuration(d time.Duration) SamplerOption {
	return func(s *sampler) error {
		s.sampleDuration = d

		return nil
	}
}

func WithSampleRatio(ratio float64) SamplerOption {
	return func(s *sampler) error {
		s.sampleRatio = ratio
		s.ratioSampler = trace.TraceIDRatioBased(ratio)

		return nil
	}
}
