package creational

import "context"

type (
	Object         interface{ any }
	ObjectExtended interface {
		Validate() error
		Init(context.Context) error
	}
)
