package trace

import "time"

type Config struct {
	AppName            string        `yaml:"app_name"`
	Ratio              float64       `yaml:"ratio" validate:"min=0,max=1"`
	SampleErrors       bool          `yaml:"sample_errors"`
	SampleLongDuration time.Duration `yaml:"sample_long_duration"`
}
