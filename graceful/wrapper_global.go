package graceful

import (
	"context"

	"github.com/agurinov/gopl/graceful/internal"
	c "github.com/agurinov/gopl/patterns/creational"
	"github.com/agurinov/gopl/run"
)

var defaultWrapper = c.Must(internal.NewWrapper())

func Run(runFn run.Fn) run.Fn {
	return defaultWrapper.Run(runFn)
}

func RunLoop(iterationFn run.Fn) run.Fn {
	return defaultWrapper.RunLoop(iterationFn)
}

func Close(closeFn run.Fn) run.Fn {
	return defaultWrapper.Close(closeFn)
}

func IsClosed(ctxs ...context.Context) bool {
	return defaultWrapper.IsClosed(ctxs...)
}
