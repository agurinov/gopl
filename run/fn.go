package run

import (
	"context"

	"github.com/agurinov/gopl/diag"
)

type (
	Fn       func(context.Context) error
	simpleFn func()
	errorFn  func() error
)

// TODO: wrapped function name is not visible
func (fn Fn) String() string {
	return diag.FunctionName(fn)
}

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
