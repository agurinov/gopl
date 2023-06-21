package gomock

import (
	"context"

	"github.com/golang/mock/gomock"
)

func IsContext() gomock.Matcher {
	return typeofMatcher[context.Context]{}
}
