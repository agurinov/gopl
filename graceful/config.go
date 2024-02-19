package graceful

import "time"

type Config struct {
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" validate:"min=200ms"`
	CloseTimeout    time.Duration `yaml:"close_timeout" validate:"min=100ms"`
}
