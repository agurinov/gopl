# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`gopl` (github.com/agurinov/gopl) is a production-grade Go platform library providing reusable infrastructure components for building cloud-native applications following 12 Factor App principles. It is a library, not a standalone application.

## Commands

```bash
# Testing
make go_test          # Unit tests with coverage
make go_itest         # Integration tests (uses Docker via dockertest)
make go_bench         # Benchmarks

# Run a single test or package
make go_test GO_PKG=./package/... TEST_USER_FLAGS='-run=^TestName'

# Code quality
make lint             # go vet + golangci-lint (30+ linters)
make fmt              # gofmt + field alignment + import fixing
make go_sec           # govulncheck security scan

# Code generation
make generate         # go generate (mocks via mockery, validators, stringers)

# Dependencies
make vendor           # go mod tidy + vendor + verify

# Build (only relevant if there's a main.go)
make go_build         # Outputs to .gura/bin/
```

Integration tests require Docker running locally.

To run only integration tests for a specific package:
```bash
make go_itest GO_PKG=./kafka/... TEST_USER_FLAGS='-run=^Test'
```

## Architecture

### Functional Options Constructor Pattern (`patterns/creational`)

Every component is constructed via generic helpers that eliminate boilerplate:

```go
var New = c.New[MyComponent, MyComponentOption]
var NewValidated = c.NewWithValidate[MyComponent, MyComponentOption]
var NewCtx = c.NewWithContext[MyComponent, MyCtxOption]
```

`Option[O]` is `func(*O) error`, `OptionWithContext[O]` is `func(context.Context, *O) error`. Components implementing `.Validate() error` get validation on construction when using `NewWithValidate`.

### Dependency Injection Convention

- `di.go` files contain the full wiring of a component's dependency graph — all `New*` calls with their options live here; wrap each constructor with `c.Must(New(...))` / `c.MustDuo(New(...))` so any construction error panics immediately (monad-style: no error threading through DI)
- `main.go` contains only bootstrap: `appcmd.Prepare(...)` + `appcmd.Start(...)`
- `graceful.Closer` collects cleanup functions during DI and runs them concurrently on shutdown (two-wave: first wave runs before second wave)
- `diag/probes.Prober` collects health-check functions during DI; after constructing a dependency call `prober.WithReadinessProbe(dep.Ping)` or `prober.WithLivenessProbe(dep.Ping)`

### Run Abstraction (`run` package)

`run.Fn = func(context.Context) error` is the core unit of execution. Use `run.Group` (stops on first error) or `run.GroupSoft` (collects all errors) for concurrent execution. `run.SimpleFn` and `run.ErrorFn` wrap functions with different signatures.

### Key Packages

| Package | Purpose |
|---|---|
| `patterns/creational` | Generic functional options constructors |
| `graceful` | Concurrent shutdown manager with wave ordering |
| `run` | Function execution, concurrency helpers, middleware |
| `appcmd` | App bootstrap: signal handling, GOMAXPROCS/GOMEMLIMIT, startup |
| `config` | Multi-source config (YAML/JSON/TOML/env/Vault) with validation |
| `diag` | Core diagnostics: metrics, tracing, status |
| `diag/log` | Zap-based structured logging |
| `diag/metrics` | Prometheus metrics |
| `diag/trace` | OpenTelemetry tracing |
| `diag/probes` | K8s-style readiness/liveness health probes |
| `http` | HTTP server with graceful shutdown, metrics, tracing |
| `grpc` | gRPC server with reflection, gRPC-Web, metrics, tracing |
| `kafka` | Kafka producer/consumer (franz-go) with metrics and hooks |
| `sql` | Database abstraction (pgx/mysql) with CQRS patterns |
| `vault` | HashiCorp Vault secrets (AppRole/UserPass auth) |
| `fsm` | Generic finite state machine with storage abstraction |
| `crontab` | Config-driven job scheduling (gocron) |
| `nopanic` | Panic recovery middleware with Telegram alerting |
| `backoff` | Retry strategies with configurable limits |
| `x` | Generic utilities: slices, maps, pointers, channels |
| `testing` | Test fixtures and assertion helpers |

### Code Generation

Generated files are suffixed `*.gen.go` and must not be edited manually:
- `*_validate.gen.go` — struct validation
- `mock_*.gen.go` — interface mocks (mockery v3)

After changing an interface or adding `//go:generate` directives, run `make generate`.

### Linting

| Rule | Value |
|---|---|
| Line length | 100 chars |
| Import order | stdlib → external → `github.com/agurinov/gopl` |
| Cyclomatic complexity | 30 |
| Cognitive complexity | 20 |
| Generated & test files | relaxed rules |

@.claude/STYLE.md
