package probes

import (
	"context"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
		closed          atomic.Bool
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

func (p *Prober) WithReadinessProbe(probes ...Probe) {
	p.readinessProbes = append(p.readinessProbes, probes...)
}

func (p *Prober) WithLivenessProbe(probes ...Probe) {
	p.livenessProbes = append(p.livenessProbes, probes...)
}

func (p *Prober) Close() {
	p.closed.Store(true)
}

func (p *Prober) Run(ctx context.Context) error {
	ticker := time.NewTicker(p.checkInterval)
	defer ticker.Stop()

	p.runAllProbes(ctx)

	p.logger.Info(
		"starting probes poller",
		zap.Stringer("check_inverval", p.checkInterval),
	)

	for {
		select {
		case <-ticker.C:
			p.runAllProbes(ctx)
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

	g, gCtx := errgroup.WithContext(checkCtx)

	for _, f := range probes {
		if f == nil {
			continue
		}

		g.Go(func() error { return f(gCtx) })
	}

	return g.Wait()
}

func (p *Prober) runAllProbes(ctx context.Context) {
	var (
		isClosed     = p.closed.Load()
		readinessErr = p.runProbes(ctx, p.readinessProbes)
		livenessErr  = p.runProbes(ctx, p.livenessProbes)
	)

	p.readiness.Store(readinessErr == nil && !isClosed)
	p.liveness.Store(livenessErr == nil)

	lvl := zapcore.DebugLevel
	if readinessErr != nil || livenessErr != nil {
		lvl = zapcore.ErrorLevel
	}

	p.logger.Log(lvl,
		"finished probes polling iteration",
		zap.NamedError("readiness_error", readinessErr),
		zap.NamedError("liveness_error", livenessErr),
	)
}
