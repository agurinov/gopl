package grpc

type Config struct {
	Port              int  `validate:"gt=1000,lt=65536"`
	ReflectionEnabled bool `yaml:"reflection_enabled"`
	MaxRequestBytes   int  `yaml:"max_request_bytes" validate:"gt=0"`
	MaxResponseBytes  int  `yaml:"max_response_bytes" validate:"gt=0"`
}
