package http

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"go.uber.org/zap"
)

func WithServerPort(port int) ServerOption {
	return func(s *Server) error {
		l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			return err
		}

		s.httpListener = l

		return nil
	}
}

func WithServerHandler(h http.Handler) ServerOption {
	return func(s *Server) error {
		if s.httpServer == nil {
			s.httpServer = new(http.Server)
		}

		s.httpServer.Handler = h

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

		s.logger = logger.Named("http.server")

		return nil
	}
}

func WithServerShutdownTimeout(t time.Duration) ServerOption {
	return func(s *Server) error {
		s.shutdownTimeout = t

		return nil
	}
}
