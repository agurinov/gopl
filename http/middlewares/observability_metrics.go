package middlewares

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	c "github.com/agurinov/gopl/patterns/creational"
)

type  observabilityMetrics struct {
}

func (m observabilityMetrics) observeServerRequest() {
}

func (m observabilityMetrics) incPanic() {
}







	onceServer.Do(func() {
		counterServerRequests = metrix.NewCounter(
			"lib_go_http_server_requests", "HTTP server requests counter",
			[]string{"method", "path", "status", "source"}, metrix.WithoutServicePrefix(),
		)
		histServerDuration = metrix.NewHistogram(
			"lib_go_http_server_duration", "HTTP server response time",
			[]string{"method", "path", "status", "source"}, metrix.WithoutServicePrefix(),
		)
	panicCounter := metrix.NewCounter(
		"lib_go_http_panic_counter", "",
		[]string{},
		metrix.WithoutServicePrefix(),
	)
	})
