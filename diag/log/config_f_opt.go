package log

func WithFormat(format string) Option {
	return func(c *Config) error {
		c.Format = format

		return nil
	}
}

func WithLevel(level string) Option {
	return func(c *Config) error {
		c.Level = level

		return nil
	}
}
