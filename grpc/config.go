package grpc

type (
	ServerConfig struct {
		Port              int  `validate:"gt=1000,lt=65536"`
		ReflectionEnabled bool `json:"reflection_enabled" yaml:"reflection_enabled"`
		DebugPayload      bool `json:"debug_payload" yaml:"debug_payload"`
		MaxRequestBytes   int  `json:"max_request_bytes" yaml:"max_request_bytes" validate:"gt=0"`
		MaxResponseBytes  int  `json:"max_response_bytes" yaml:"max_response_bytes" validate:"gt=0"`
	}
	ClientConfig struct {
		Addr      string `validate:"required"`
		AuthToken string
	}
)
