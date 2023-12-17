package creational

import "context"

type (
	Option[O Object]      func(*O) error
	optionAlias[O Object] interface {
		~func(*O) error
	}
)

type (
	OptionWithContext[O Object]      func(context.Context, *O) error
	optionWithContextAlias[O Object] interface {
		~func(context.Context, *O) error
	}
)
