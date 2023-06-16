package gomock

import (
	"fmt"

	"github.com/golang/mock/gomock"
)

type equalable[T any] interface {
	Equal(T) bool
}

type equalMethodMatcher[T any] struct {
	want T
}

func (e equalMethodMatcher[T]) Matches(got any) bool {
	var (
		gotTyped, gotOk   = got.(T)
		wantTyped, wantOk = any(e.want).(equalable[T])
	)

	if gotOk && wantOk {
		return wantTyped.Equal(gotTyped)
	}

	return gomock.Eq(e.want).Matches(got)
}

func (e equalMethodMatcher[T]) String() string {
	return fmt.Sprintf("is equal to %v (%T) via .Equal() method", e.want, e.want)
}

func Eq[T any](want T) gomock.Matcher {
	return equalMethodMatcher[T]{
		want: want,
	}
}
