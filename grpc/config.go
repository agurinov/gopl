package grpc

type Config struct {
	Port              int `validate:"gt=1000,lt=65536"`
	ReflectionEnabled bool
	MaxRequestBytes   int `validate:"gt=0"`
	MaxResponseBytes  int `validate:"gt=0"`
}
