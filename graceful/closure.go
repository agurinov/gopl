package graceful

import (
	"context"

	"github.com/agurinov/gopl/diag"
)

type (
	Closure       func(context.Context) error
	simpleClosure func()
	errorClosure  func() error
)

// TODO: wrapped function name is not visible
func (c Closure) String() string {
	return diag.FunctionName(c)
}

func (c simpleClosure) f(context.Context) error {
	c()

	return nil
}

func (c errorClosure) f(context.Context) error {
	return c()
}

func SimpleClosure(fn func()) Closure {
	return simpleClosure(fn).f
}

func ErrorClosure(fn func() error) Closure {
	return errorClosure(fn).f
}
