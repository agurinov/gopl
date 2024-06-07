package middlewares

import (
	"fmt"
	"net/http"

	"github.com/agurinov/gopl/diag/trace"
)

func Trace(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if r.Header != nil {
			ctx = trace.TraceparentToContext(ctx, r.Header.Get("traceparent"))
		}

		recorder := &statusRecorder{
			ResponseWriter: w,
			Status:         http.StatusOK,
		}

		ctx, span := trace.StartSpan(ctx, "http.Router")
		defer func() {
			switch recorder.Status {
			case http.StatusInternalServerError:
				trace.RegisterError(span, fmt.Errorf("possible panic"))
			case http.StatusNotFound, http.StatusMethodNotAllowed:
				return
			}

			span.End()
		}()

		next.ServeHTTP(recorder, r.WithContext(ctx))
	})
}
