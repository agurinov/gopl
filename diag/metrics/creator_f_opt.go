package metrics

func WithoutServicePrefix() Option {
	return func(c *creator) error {
		c.noServicePrefix = true

		return nil
	}
}

func WithBuckets(buckets []float64) Option {
	return func(c *creator) error {
		c.buckets = buckets

		return nil
	}
}
