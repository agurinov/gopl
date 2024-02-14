// Code generated: TODO

package http

import (
	"net"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	"go.uber.org/zap"
)

func (obj Server) Validate() error {
	s := struct {
		HttpListener    net.Listener  `validate:"required"`
		Logger          *zap.Logger   `validate:"required"`
		HttpServer      *http.Server  `validate:"required"`
		Name            string        `validate:"required"`
		ShutdownTimeout time.Duration `validate:"required"`
	}{
		Name:            obj.name,
		HttpListener:    obj.httpListener,
		HttpServer:      obj.httpServer,
		Logger:          obj.logger,
		ShutdownTimeout: obj.shutdownTimeout,
	}

	v := validator.New()
	if err := v.RegisterValidation("notblank", validators.NotBlank); err != nil {
		return err
	}

	if err := v.Struct(s); err != nil {
		return err
	}

	return nil
}
