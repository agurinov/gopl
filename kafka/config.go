package kafka

import "fmt"

type Config struct {
	EventPosition
	BatchSize uint
}

func (c Config) Validate() error {
	if c.Topic == "" {
		return fmt.Errorf(
			"%w: topic name must be present",
			ErrInvalidConfig,
		)
	}

	if c.BatchSize == 0 {
		return fmt.Errorf(
			"%w: batch size must be gt 0",
			ErrInvalidConfig,
		)
	}

	return nil
}
