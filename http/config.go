package http

type Config struct {
	Name string
	Port int `validate:"gt=1000,lt=65536"`
}

func (c Config) NewServer(opts ...ServerOption) (Server, error) {
	defaults := []ServerOption{
		WithServerName(c.Name),
		WithServerPort(c.Port),
	}

	opts = append(defaults, opts...)

	return NewServer(opts...)
}
