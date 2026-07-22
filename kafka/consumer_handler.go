package kafka

import (
	"context"
)

type (
	Handler      = func(context.Context, Record) error
	HandlerBatch = func(context.Context, []Record) error
)
