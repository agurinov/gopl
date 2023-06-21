package creational

import "context"

type (
	Object interface{ any }

	Option[O Object]     func(*O) error
	OptionFunc[O Object] interface {
		~func(*O) error
	}

	OptionWithContext[O Object]     func(context.Context, *O) error
	OptionFuncWithContext[O Object] interface {
		~func(context.Context, *O) error
	}
)
