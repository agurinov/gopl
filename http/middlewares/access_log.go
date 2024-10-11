package middlewares

import (
	"net/http"
	"time"

	"github.com/dustin/go-humanize"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func AccessLog(logger *zap.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			recorder := &statusRecorder{
				ResponseWriter: w,
				Status:         http.StatusOK,
			}

			startTime := time.Now()

			next.ServeHTTP(recorder, r)

			elapsedTime := time.Since(startTime)

			logLvl := zapcore.InfoLevel
			if recorder.Status >= 500 && recorder.Status <= 599 {
				logLvl = zapcore.ErrorLevel
			}

			reqContentLen := "unknown"
			if r.ContentLength >= 0 {
				// skip integer overflow conversion int64 -> uint64
				//nolint:gosec
				reqContentLen = humanize.Bytes(uint64(r.ContentLength))
			}

			logger.Log(logLvl,
				"http served request",
				zap.Int("status_code", recorder.Status),
				zap.String("remote_addr", r.RemoteAddr),
				zap.String("http_method", r.Method),
				zap.String("request_uri", r.RequestURI),
				zap.String("request_content_length", reqContentLen),
				zap.Stringer("elapsed_time", elapsedTime),
			)
		})
	}
}
