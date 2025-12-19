package graceful

import (
	"sync"

	"github.com/agurinov/gopl/diag/log"
	"github.com/agurinov/gopl/run"
)

var (
	defaultWrapper     Wrapper
	defaultWrapperOnce sync.Once
)

func initDefaultWrapper() {
	logger := log.MustNewZapSystem()

	w, err := NewWrapper(
		WithWrapperLogger(logger),
	)
	if err != nil {
		panic(err)
	}

	defaultWrapper = w
}

func WrapRun(f run.Closure) run.Closure {
	defaultWrapperOnce.Do(initDefaultWrapper)

	return defaultWrapper.WrapRun(f)
}

func WrapClose(f run.Closure) run.Closure {
	defaultWrapperOnce.Do(initDefaultWrapper)

	return defaultWrapper.WrapClose(f)
}
