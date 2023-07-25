package gomock

import (
	"fmt"

	"github.com/golang/mock/gomock"
)

type typeofMatcher[T any] struct{}

func (typeofMatcher[T]) Matches(got any) bool {
	_, gotOk := got.(T)

	return gotOk
}

func (typeofMatcher[T]) String() string {
	var t T

	return fmt.Sprintf("is typeof %T", t)
}

func TypeOf[T any]() gomock.Matcher {
	return typeofMatcher[T]{}
}
