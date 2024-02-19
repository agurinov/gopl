package http

type Config struct {
	Port int `validate:"gt=1000,lt=65536"`
}
