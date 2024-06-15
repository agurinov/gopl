package trace

import (
	"fmt"
	"time"

	"go.opentelemetry.io/otel/sdk/trace"

	c "github.com/agurinov/gopl/patterns/creational"
)

type (
	sampler struct {
		ratioSampler   trace.Sampler
		sampleRatio    float64
		sampleDuration time.Duration
		sampleErrors   bool
	}
	SamplerOption = c.Option[sampler]
)

func NewSampler(opts ...SamplerOption) (trace.Sampler, error) {
	return c.NewWithValidate[sampler, SamplerOption](opts...)
}

func (s sampler) Description() string {
	return fmt.Sprintf(
		"sampler=gopl errors=%t duration=%s ratio=%f",
		s.sampleErrors,
		s.sampleDuration,
		s.sampleRatio,
	)
}

func (s sampler) ShouldSample(params trace.SamplingParameters) trace.SamplingResult {
	if s.sampleErrors && isErrorSpan(params) {
		return trace.SamplingResult{
			Decision: trace.RecordAndSample,
		}
	}

	if s.sampleDuration != 0 && isLongSpan(params) {
		return trace.SamplingResult{
			Decision: trace.RecordAndSample,
		}
	}

	return s.ratioSampler.ShouldSample(params)
}

func isErrorSpan(params trace.SamplingParameters) bool {
	for i := range params.Attributes {
		if params.Attributes[i].Key == "error" && params.Attributes[i].Value.AsBool() {
			return true
		}
	}

	return false
}

func isLongSpan(trace.SamplingParameters) bool {
	return false
}
