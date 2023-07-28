package creational

import "context"

type (
	Option[O Object]            func(*O) error
	OptionWithContext[O Object] func(context.Context, *O) error

	OptionFunc[O Object] interface {
		~func(*O) error
	}
	OptionFuncWithContext[O Object] interface {
		~func(context.Context, *O) error
	}
)
