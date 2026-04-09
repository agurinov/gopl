# graceful

Concurrent shutdown manager designed for the **DI stage** (dependency injection).
It collects cleanup steps and runs them **concurrently** when your application context is cancelled.

- Collect `func()`, `func(context.Context)`, `func(context.Context) error` functions during DI wiring.
- Trigger shutdown via `closer.WaitForShutdown(context.Context)`.
- All steps run **in parallel** (via `errgroup`). Errors are aggregated with `errors.Join`.
- Functional options:
  - `WithLogger(*zap.Logger)` — required via validation.
  - `WithTimeout(time.Duration)` — required via validation (must be > 0).

> This package assumes closers are independent or resilient to parallel shutdown.
> If you have strict dependency order, wrap those closers in a single step (or serialize manually).

https://victoriametrics.com/blog/go-graceful-shutdown/

---

## Usage

```go
// di.go

import (
	c "github.com/agurinov/gopl/patterns/creational"
	"github.com/agurinov/gopl/graceful"
	"github.com/agurinov/gopl/run"
)

type diContainer struct {
	closer *graceful.Closer
}

closer := c.Must(
	graceful.NewCloser(
		graceful.WithLogger(logger),
		graceful.WithTimeout(30*time.Second),
	),
)

svc1 := c.Must(
	// ...
)
closer.AddCloser(
	svc1.Close,
	graceful.InFirstWave(),
)

svc2 := c.Must(
	// ...
)
closer.AddCloser(
	run.SimpleFn(svc2.Close),
)

svc3 := c.Must(
	// ...
)
closer.AddCloser(
	run.ErrorFn(svc3.Close),
)

return diContainer{
	closer: closer,
}
```

```go
// main.go

import cmd "github.com/agurinov/gopl/appcmd"

const cmdName = "api"

func main() {
	ctx, stop, logger, err := cmd.Prepare(cmdName)
	if err != nil {
		logger.Fatal("can't init cmd", zap.Error(err))
	}
	defer stop()

	di := newDI(ctx, cfg)

	cmd.Start(
		ctx,
		logger,
		di.closer.WaitForShutdown,
	)
}
```
