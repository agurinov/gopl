package probes

import (
	"context"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	c "github.com/agurinov/gopl/patterns/creational"
	"github.com/agurinov/gopl/run"
)

type (
	ProbeType string
	Prober    struct {
		logger          *zap.Logger
		readinessProbes []run.Fn
		livenessProbes  []run.Fn
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
	ProbeTypeReadiness ProbeType = "readiness"
	ProbeTypeLiveness  ProbeType = "liveness"
)

func New(opts ...Option) (*Prober, error) {
	return c.ConstructWithValidate(new(Prober), opts...)
}

func (p *Prober) Startup() bool   { return p.startup.Load() }
func (p *Prober) Readiness() bool { return p.readiness.Load() }
func (p *Prober) Liveness() bool  { return p.liveness.Load() }

func (p *Prober) SetStartup(startup bool) {
	p.startup.Store(startup)
}

func (p *Prober) WithReadinessProbe(probes ...run.Fn) {
	p.readinessProbes = append(p.readinessProbes, probes...)
}

func (p *Prober) WithLivenessProbe(probes ...run.Fn) {
	p.livenessProbes = append(p.livenessProbes, probes...)
}

func (p *Prober) Close() {
	p.closed.Store(true)

	p.logger.Info("prober closed; readiness always false")
}

func (p *Prober) Run(ctx context.Context) error {
	p.runAllProbes(ctx)

	ticker := time.NewTicker(p.checkInterval)
	defer ticker.Stop()

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
	probes []run.Fn,
) error {
	if len(probes) == 0 {
		return nil
	}

	checkCtx, checkCancel := context.WithTimeout(ctx, p.checkTimeout)
	defer checkCancel()

	return run.Group(checkCtx, probes...)
}

func (p *Prober) runAllProbes(ctx context.Context) {
	var (
		isClosed     = p.closed.Load()
		readinessErr = p.runProbes(ctx, p.readinessProbes)
		livenessErr  = p.runProbes(ctx, p.livenessProbes)
	)

	var (
		startup   = p.Startup()
		readiness = readinessErr == nil && !isClosed
		liveness  = livenessErr == nil
	)

	p.readiness.Store(readiness)
	p.liveness.Store(liveness)

	lvl := zapcore.DebugLevel

	switch {
	case readinessErr != nil || livenessErr != nil:
		lvl = zapcore.ErrorLevel
	case isClosed:
		lvl = zapcore.InfoLevel
	}

	p.logger.Log(lvl,
		"finished probes polling iteration",
		zap.Bool("startup", startup),
		zap.Bool("readiness", readiness),
		zap.Bool("liveness", liveness),
	)
}
