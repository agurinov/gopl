package grpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	c "github.com/agurinov/gopl/patterns/creational"
)

type (
	Server struct {
		grpcListener         net.Listener
		logger               *zap.Logger
		grpcServer           *grpc.Server
		grpcServices         map[*grpc.ServiceDesc]any
		name                 string
		shutdownTimeout      time.Duration
		grpcServerReflection bool
	}
	ServerOption c.Option[Server]
)

var NewServer = c.NewWithValidate[Server, ServerOption]

func (s Server) GRPC() *grpc.Server {
	serviceInfo := s.grpcServer.GetServiceInfo()

	if s.grpcServerReflection {
		serviceName := "grpc.reflection.v1alpha.ServerReflection"
		if _, registered := serviceInfo[serviceName]; !registered {
			reflection.Register(s.grpcServer)
		}
	}

	for desc, impl := range s.grpcServices {
		if _, registered := serviceInfo[desc.ServiceName]; registered {
			continue
		}

		s.grpcServer.RegisterService(desc, impl)
	}

	return s.grpcServer
}

func (s Server) GRPCWeb() http.Handler {
	return grpcweb.WrapServer(s.GRPC())
}

func (s Server) ListenAndServe(context.Context) error {
	s.logger.Info(
		"starting grpc server",
		zap.String("server_name", s.name),
		zap.Stringer("server_address", s.grpcListener.Addr()),
	)

	switch err := s.GRPC().Serve(s.grpcListener); {
	case err == nil:
	case errors.Is(err, grpc.ErrServerStopped):
	default:
		return fmt.Errorf("server %q: can't listen: %w", s.name, err)
	}

	return nil
}

func (s Server) Stop() {
	s.logger.Info(
		"shutting down grpc server",
		zap.String("server_name", s.name),
		zap.Stringer("server_address", s.grpcListener.Addr()),
	)

	s.GRPC().GracefulStop()
}

// Deprecated: use closer.AddCloser(server.Stop) instead
func (s Server) WaitForShutdown(ctx context.Context) error {
	<-ctx.Done()

	s.Stop()

	return nil
}
