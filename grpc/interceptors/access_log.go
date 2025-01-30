package interceptors

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/agurinov/gopl/diag/log"
)

//nolint:revive
func LoggerUnaryServer(
	logger *zap.Logger,
	debugPayload bool,
) grpc.UnaryServerInterceptor {
	loggableEvents := []logging.LoggableEvent{
		logging.StartCall,
		logging.FinishCall,
	}

	if debugPayload {
		loggableEvents = append(loggableEvents,
			logging.PayloadReceived,
			logging.PayloadSent,
		)
	}

	return logging.UnaryServerInterceptor(
		log.GRPC(logger),
		logging.WithDurationField(
			logging.DurationToDurationField,
		),
		logging.WithLogOnEvents(loggableEvents...),
		logging.WithDisableLoggingFields(
			"protocol",
			"peer.address",
			"grpc.component",
			"grpc.start_time",
			"grpc.method_type",
		),
		logging.WithLevels(serverCodeToLevel),
	)
}

func serverCodeToLevel(code codes.Code) logging.Level {
	switch code {
	case codes.OK,
		codes.NotFound,
		codes.Canceled,
		codes.AlreadyExists,
		codes.Unauthenticated:
		return logging.LevelInfo

	case codes.DeadlineExceeded,
		codes.PermissionDenied,
		codes.ResourceExhausted,
		codes.FailedPrecondition,
		codes.Aborted,
		codes.OutOfRange,
		codes.Unavailable,
		codes.InvalidArgument:
		return logging.LevelWarn

	case codes.Unknown,
		codes.Unimplemented,
		codes.Internal,
		codes.DataLoss:
		return logging.LevelError

	default:
		return logging.LevelError
	}
}
