package trace

import (
	"context"
	"time"
)

type Config struct {
	AppName            string        `json:"app_name" yaml:"app_name" validate:"required"`
	Ratio              float64       `json:"ratio" yaml:"ratio" validate:"gte=0.01,lte=1"`
	SampleErrors       bool          `json:"sample_errors" yaml:"sample_errors"`
	SampleLongDuration time.Duration `json:"sample_long_duration" yaml:"sample_long_duration"`
}

func (c Config) Init(
	ctx context.Context,
	opts ...IniterOption,
) error {
	defaults := []IniterOption{
		WithCmdName(c.AppName),
		WithSamplerOptions(
			WithSampleError(c.SampleErrors),
			WithSampleDuration(c.SampleLongDuration),
			WithSampleRatio(c.Ratio),
		),
	}

	opts = append(defaults, opts...)

	return Init(ctx, opts...)
}
