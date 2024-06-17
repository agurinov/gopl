package trace

import "time"

type Config struct {
	AppName            string        `yaml:"app_name" validate:"required"`
	Ratio              float64       `yaml:"ratio" validate:"gte=0.01,lte=1"`
	SampleErrors       bool          `yaml:"sample_errors"`
	SampleLongDuration time.Duration `yaml:"sample_long_duration"`
}
