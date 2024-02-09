package middlewares

import (
	"fmt"
	"net/http"
	"net/url"
	"runtime/debug"
	"strconv"
	"time"

	c "github.com/agurinov/gopl/patterns/creational"
)

type (
	observability struct {
		logger     *zap.logger
		metrics    observabilityMetrics
		leadingLog bool
	}
	ObservabilityOption c.Option[observability]
)

func NewObservability(opts ...ObservabilityOption) (Middleware, error) {
	m, err := c.NewWithValidate[observability, ObservabilityOption](opts...)
	if err != nil {
		return nil, err
	}

	return m.Middleware, nil
}

func (o observability) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if traceparent := r.Header.Get("traceparent"); traceparent != "" {
			ctx = trace.TraceparentToContext(ctx, traceparent)
		}

		recorder := &statusRecorder{
			ResponseWriter: w,
			Status:         http.StatusOK,
		}

		ctx, span := trace.StartSpan(ctx, "http.Observability")
		defer func() {
			switch recorder.Status {
			case http.StatusNotFound, http.StatusMethodNotAllowed:
			default:
				span.End()
			}
		}()

		defer func() {
			if r := recover(); r != nil {
				var (
					err   = fmt.Errorf("%s", r)
					stack = string(debug.Stack())
				)

				o.logger.Error(
					"panic recovered",
					zap.String("stack", stack),
					zap.Error(err),
				)
				trace.RegisterError(span, err)
				o.metrics.incPanic()

				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		// Leading access log

		startTime := time.Now()
		next.ServeHTTP(recorder, r.WithContext(ctx))
		elapsedTime := time.Since(startTime)

		res := strconv.Itoa(rw.status)
		source := parseSource(r)

		if rw.status != http.StatusNotFound {
			var cfg MetricsConfig
			cfg.FormatPath = func(r *http.Request) string { return r.URL.Path }

			for _, opt := range opts {
				opt(&cfg)
			}

			path = cfg.FormatPath(r)
			u, err := url.Parse(path)
			if err == nil && u.Path != "" {
				path = u.Path
			}
		} else {
			path = "http404"
		}

		m.metrics.observeServerRequest(r.Method, path, res, source)

		logger.Info("http served request",
			zap.Int("status_code", recorder.Status),
			zap.String("remote_addr", r.RemoteAddr),
			zap.String("http_method", r.Method),
			zap.String("request_uri", r.RequestURI),
			zap.String("content_length", humanize.Bytes(uint64(r.ContentLength))),
			zap.Stringer("elapsed_time", elapsedTime),
		)
	})

}
