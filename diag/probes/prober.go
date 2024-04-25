package probes

import (
	"context"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	c "github.com/agurinov/gopl/patterns/creational"
)

type (
	Probe     = func(context.Context) error
	ProbeType string
	Prober    struct {
		logger          *zap.Logger
		readinessProbes []Probe
		livenessProbes  []Probe
		checkInterval   time.Duration
		checkTimeout    time.Duration
		startup         atomic.Bool
		readiness       atomic.Bool
		liveness        atomic.Bool
	}
	Option = c.Option[*Prober]
)

const (
	ProbeTypeStartup   ProbeType = "startup"
	ProbeTypeReadiness ProbeType = "readiness"
	ProbeTypeLiveness  ProbeType = "liveness"
)

func New(opts ...Option) (*Prober, error) {
	return c.ConstructWithValidate(new(Prober), opts...)
}

func (p *Prober) SetStartup(startup bool) { p.startup.Store(startup) }
func (p *Prober) Startup() bool           { return p.startup.Load() }
func (p *Prober) Readiness() bool         { return p.readiness.Load() }
func (p *Prober) Liveness() bool          { return p.liveness.Load() }

func (p *Prober) Run(ctx context.Context) error {
	ticker := time.NewTicker(p.checkInterval)
	defer ticker.Stop()

	p.logger.Info(
		"starting probes poller",
		zap.Duration("check_inverval", p.checkInterval),
	)

	for {
		select {
		case <-ticker.C:
			readinessErr := p.runProbes(ctx, p.readinessProbes)
			p.readiness.Store(readinessErr == nil)

			livenessErr := p.runProbes(ctx, p.livenessProbes)
			p.liveness.Store(livenessErr == nil)

			p.logger.Debug(
				"finished probes polling iteration",
				zap.NamedError("readiness_error", readinessErr),
				zap.NamedError("liveness_error", livenessErr),
			)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (p *Prober) runProbes(
	ctx context.Context,
	probes []func(context.Context) error,
) error {
	if len(probes) == 0 {
		return nil
	}

	checkCtx, checkCancel := context.WithTimeout(ctx, p.checkTimeout)
	defer checkCancel()

	g, ctx := errgroup.WithContext(checkCtx)

	for _, f := range probes {
		if f == nil {
			continue
		}

		f := f

		g.Go(func() error { return f(ctx) })
	}

	return g.Wait()
}
