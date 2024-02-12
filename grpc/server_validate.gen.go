package grpc

import (
	"net"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func (obj Server) Validate() error {
	s := struct {
		GrpcListener    net.Listener              `validate:"required"`
		Logger          *zap.Logger               `validate:"required"`
		GrpcServer      *grpc.Server              `validate:"required"`
		GrpcServices    map[*grpc.ServiceDesc]any `validate:"gt=0,dive,keys,required,endkeys,required"`
		Name            string                    `validate:"required"`
		ShutdownTimeout time.Duration             `validate:"required"`
	}{
		Name:            obj.name,
		Logger:          obj.logger,
		GrpcListener:    obj.grpcListener,
		GrpcServer:      obj.grpcServer,
		GrpcServices:    obj.grpcServices,
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
