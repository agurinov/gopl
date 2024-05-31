package sql

func WithRW(db DB) ClientOption {
	return func(c *cqrs) error {
		c.rw = db

		return nil
	}
}

func WithRO(db DB) ClientOption {
	return func(c *cqrs) error {
		c.ro = db

		return nil
	}
}
