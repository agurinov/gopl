# HTTP handlers

## Table of contents:
- [Debug handler](#debug)
- [Basic handler](#basic)
- [Static handler](#static)

### Debug handler

```go
// di.go

import (
	c "github.com/agurinov/gopl/patterns/creational"
	"github.com/agurinov/gopl/http/handlers"
)

type diContainer struct {
	debugServer http.Server
}

// Zap atomic level to be changed on the fly.
logger, lvl := c.MustDuo(
	cfg.Logger.New(),
)

nopanicHandler := c.Must(
	nopanic.NewHandler(),
)

// Orchestrator prober.
prober := c.Must(
	cfg.Probes.NewProber(),
)

debugHandler := c.Must(
	handlers.NewDebug(
		handlers.WithDebugAtomicLevel(lvl),
		handlers.WithDebugLogger(logger),
		handlers.WithDebugCustomMiddlewares(
			nopanicHandler.Middleware(),
		),
		handlers.WithDebugProber(prober),
	),
)

// Server to serve all debugHandler endpoints.
debugServer := c.Must(
	cfg.HTTP.NewServer(),
)

return diContainer{
	debugServer: debugServer,
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
		di.debugServer.ListenAndServe,
	)
}
```

### Basic handler

### Static handler
