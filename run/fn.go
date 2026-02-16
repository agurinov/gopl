package run

import (
	"context"
)

type (
	Fn       = func(context.Context) error
	simpleFn func()
	errorFn  func() error
)

func (fn simpleFn) f(context.Context) error {
	fn()

	return nil
}

func (fn errorFn) f(context.Context) error {
	return fn()
}

func SimpleFn(fn func()) Fn {
	return simpleFn(fn).f
}

func ErrorFn(fn func() error) Fn {
	return errorFn(fn).f
}
