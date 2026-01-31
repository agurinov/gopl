package grpc

import (
	"errors"
	"fmt"
	"net"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

func WithServerPort(port int) ServerOption {
	return func(s *Server) error {
		l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			return err
		}

		s.grpcListener = l

		return nil
	}
}

func WithServerReflection(enabled bool) ServerOption {
	return func(s *Server) error {
		s.grpcServerReflection = enabled

		return nil
	}
}

func WithServerName(name string) ServerOption {
	return func(s *Server) error {
		s.name = name

		return nil
	}
}

func WithServerLogger(logger *zap.Logger) ServerOption {
	return func(s *Server) error {
		if logger == nil {
			return nil
		}

		s.logger = logger.Named("grpc.server")

		grpclog.SetLoggerV2(
			zapgrpc.NewLogger(
				logger.Named("grpc.transport"),
			),
		)

		return nil
	}
}

// Deprecated: Use closer centralized mechanics instead.
func WithServerShutdownTimeout(shutdownTimeout time.Duration) ServerOption {
	return func(s *Server) error {
		s.shutdownTimeout = shutdownTimeout

		return nil
	}
}

func WithService(desc *grpc.ServiceDesc, impl any) ServerOption {
	return func(s *Server) error {
		if s.grpcServices == nil {
			s.grpcServices = make(map[*grpc.ServiceDesc]any)
		}

		if _, exists := s.grpcServices[desc]; exists {
			return fmt.Errorf("service %q already defined", desc.ServiceName)
		}

		s.grpcServices[desc] = impl

		return nil
	}
}

func WithServerOptions(opts ...grpc.ServerOption) ServerOption {
	return func(s *Server) error {
		if s.grpcServer != nil {
			return errors.New("server options already binded")
		}

		s.grpcServer = grpc.NewServer(opts...)

		return nil
	}
}
