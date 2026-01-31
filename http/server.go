package http

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"go.uber.org/zap"

	c "github.com/agurinov/gopl/patterns/creational"
)

type (
	Server struct {
		httpListener    net.Listener
		logger          *zap.Logger
		httpServer      *http.Server
		name            string
		shutdownTimeout time.Duration
	}
	ServerOption c.Option[Server]
)

var NewServer = c.NewWithValidate[Server, ServerOption]

func (s Server) ListenAndServe(context.Context) error {
	s.logger.Info(
		"starting http server",
		zap.String("server_name", s.name),
		zap.Stringer("server_address", s.httpListener.Addr()),
	)

	s.httpServer.RegisterOnShutdown(func() {
		s.logger.Info(
			"shutting down server",
			zap.String("server_name", s.name),
			zap.Stringer("server_address", s.httpListener.Addr()),
		)
	})

	switch err := s.httpServer.Serve(s.httpListener); {
	case err == nil:
	case errors.Is(err, http.ErrServerClosed):
	default:
		return fmt.Errorf("server %q: can't listen: %w", s.name, err)
	}

	return nil
}

func (s Server) Shutdown(ctx context.Context) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server %q: can't shutdown: %w", s.name, err)
	}

	return nil
}

// Deprecated: use closer.AddCloser(server.Shutdown) instead
func (s Server) WaitForShutdown(ctx context.Context) error {
	<-ctx.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(
		context.Background(),
		s.shutdownTimeout,
	)
	defer shutdownCancel()

	//nolint:contextcheck
	return s.Shutdown(shutdownCtx)
}
