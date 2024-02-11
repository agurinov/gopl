package log

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.uber.org/zap"
)

func GRPC(logger *zap.Logger) logging.Logger {
	sl := logger.Named("grpc.handler").Sugar()

	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		var logfn func(string, ...any)

		switch lvl {
		case logging.LevelDebug:
			logfn = sl.Debugw
		case logging.LevelInfo:
			logfn = sl.Infow
		case logging.LevelWarn:
			logfn = sl.Warnw
		default:
			logfn = sl.Errorw
		}

		logfn(msg, fields...)
	})
}
