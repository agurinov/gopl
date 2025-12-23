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

func Run(fn run.Fn) run.Fn {
	defaultWrapperOnce.Do(initDefaultWrapper)

	return defaultWrapper.Run(fn)
}

func Close(fn run.Fn) run.Fn {
	defaultWrapperOnce.Do(initDefaultWrapper)

	return defaultWrapper.Close(fn)
}
