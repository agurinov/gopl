package probes

import "time"

type Config struct {
	CheckInterval time.Duration `yaml:"check_interval" validate:"min=1s"`
	CheckTimeout  time.Duration `yaml:"check_timeout" validate:"min=1s"`
}
