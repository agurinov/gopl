package trace

/*
import (
	"fmt"
	"os"

	"go.opentelemetry.io/otel/exporters/jaeger"
)

func NewJaegerExporter(host, port string) (*jaeger.Exporter, error) {
	exporter, err := jaeger.New(
		jaeger.WithAgentEndpoint(
			jaeger.WithAgentHost(host),
			jaeger.WithAgentPort(port),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("can't init Jaeger exporter `%s:%s`: %w", host, port, err)
	}
	return exporter, nil
}

func NewJaegerExporterFromEnv() (*jaeger.Exporter, error) {
	jaegerHost := os.Getenv("JAEGER_AGENT_HOST")
	if jaegerHost == "" {
		jaegerHost = "localhost"
	}
	jaegerPort := os.Getenv("JAEGER_AGENT_PORT")
	if jaegerPort == "" {
		jaegerPort = "6831"
	}
	return NewJaegerExporter(jaegerHost, jaegerPort)
}
*/
