package interceptors

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/agurinov/gopl/diag/log"
)

func LoggerUnaryServer(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return logging.UnaryServerInterceptor(
		log.GRPC(logger),
		logging.WithLogOnEvents(
			logging.StartCall,
			logging.FinishCall,
		),
		logging.WithDisableLoggingFields("protocol", "grpc.component", "grpc.start_time"),
	)
}
