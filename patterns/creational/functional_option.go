package creational

import "context"

type (
	Option[O Object]     func(*O) error
	OptionFunc[O Object] interface {
		~func(*O) error
	}
)

type (
	OptionWithContext[O Object]     func(context.Context, *O) error
	OptionFuncWithContext[O Object] interface {
		~func(context.Context, *O) error
	}
)
