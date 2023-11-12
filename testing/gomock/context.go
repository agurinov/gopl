package gomock

import (
	"context"

	"go.uber.org/mock/gomock"
)

func IsContext() gomock.Matcher {
	return typeofMatcher[context.Context]{}
}
