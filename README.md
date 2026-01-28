# Golang platform library

![GitHub license](https://img.shields.io/github/license/agurinov/gopl)
![GitHub repo size](https://img.shields.io/github/repo-size/agurinov/gopl)
[![Latest Tag](https://img.shields.io/github/v/tag/agurinov/gopl?label=latest&logo=git)](https://github.com/agurinov/gopl/tags)

[![Go Reference](https://pkg.go.dev/badge/github.com/agurinov/gopl.svg)](https://pkg.go.dev/github.com/agurinov/gopl)
[![Go Report Card](https://goreportcard.com/badge/github.com/agurinov/gopl)](https://goreportcard.com/report/github.com/agurinov/gopl)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/agurinov/gopl)

![CI](https://github.com/agurinov/gopl/actions/workflows/golang-library.yaml/badge.svg)
![golangci-lint](https://img.shields.io/badge/linter-golangci--lint-blue)

## Table of contents:
- Patterns (`patterns`)
	- [Creational](./patterns/creational/README.md)
- Diagnotics (`diag`)
	- [probes](./diag/probes/README.md)
- [Graceful](./graceful/README.md)

## Architectural Approach

The majority of packages in this repository implement:

- **12 Factor App** best practices
- Ready-to-use integrations with infrastructure (logging, databases, observability, etc.)

Examples usually demonstrate two key layers of application assembly:
```go
// main.go

// The executable entry point that runs inside a container.
// Minimal logic — just bootstrap and start.
```

 ```go
// di.go

// Explicit construction of the full dependency graph.
// All services, repositories, clients and middleware are wired together in a single, readable place — making dependency flow transparent and refactor-friendly.
```
