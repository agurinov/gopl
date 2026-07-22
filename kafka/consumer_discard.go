package kafka

import (
	"context"
)

type (
	RecordDiscarder interface {
		Discard(context.Context, ...Record) error
	}
)
