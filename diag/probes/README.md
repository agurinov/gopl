# probes

## Usage

```go
// di.go

import (
	"github.com/agurinov/gopl/diag/probes"
	c "github.com/agurinov/gopl/patterns/creational"
)

type diContainer struct {
	prober *probes.Prober
}

prober := c.Must(
	probes.New(
		probes.WithLogger(logger),
		probes.WithCheckInterval(5*time.Second),
		probes.WithCheckTimeout(time.Second),
	),
)

// Good idea to stop traffic from orchestrator via readiness probe.
closer.AddCloser(
	run.SimpleFn(prober.Close),
	graceful.InFirstWave(),
)

sqlx := c.Must(
	// ...
)

// db is a critical component. 
// No ping -> no readiness probe -> no traffic from orchestrator.
prober.WithReadinessProbe(sqlx.Ping)

// After all dependency tree is ready we talk to orchestrator that app is started.
prober.SetStartup(true)

return diContainer{
	prober: prober,
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
		di.prober.Run,
	)
}
```
