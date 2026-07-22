package metrics

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

func registerAs[T any](
	cr creator,
	vec prometheus.Collector,
) T {
	registered := cr.register(vec)

	typedVec, ok := registered.(T)
	if !ok {
		var zero T

		panic(
			fmt.Sprintf(
				"metrics: registered collector is not a %T",
				zero,
			),
		)
	}

	return typedVec
}
