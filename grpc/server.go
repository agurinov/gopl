package grpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

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

func (s Server) ListenAndServe(_ context.Context) error {
	s.logger.Info(
		"starting grpc server",
		zap.String("server_name", s.name),
		zap.Stringer("server_address", s.grpcListener.Addr()),
	)

	if s.grpcServerReflection {
		reflection.Register(s.grpcServer)
	}

	for desc, impl := range s.grpcServices {
		s.grpcServer.RegisterService(desc, impl)
	}

	switch err := s.grpcServer.Serve(s.grpcListener); {
	case err == nil:
	case errors.Is(err, grpc.ErrServerStopped):
	default:
		return fmt.Errorf("server %q: can't listen: %w", s.name, err)
	}

	return nil
}

func (s Server) WaitForShutdown(ctx context.Context) error {
	<-ctx.Done()

	s.logger.Debug(
		"shutting down grpc server",
		zap.String("server_name", s.name),
		zap.Stringer("server_address", s.grpcListener.Addr()),
	)

	s.grpcServer.GracefulStop()

	return nil
}
