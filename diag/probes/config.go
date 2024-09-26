package probes

import "time"

type Config struct {
	CheckInterval time.Duration `json:"check_interval" yaml:"check_interval" validate:"min=1s"`
	CheckTimeout  time.Duration `json:"check_timeout" yaml:"check_timeout" validate:"min=1s"`
}
