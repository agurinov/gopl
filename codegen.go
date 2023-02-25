//go:build neverbuild

package main

// Define generate aliases (long shared commands)
//go:generate -command stringer go run ./vendor/golang.org/x/tools/cmd/stringer
//go:generate -command mockgen go run ./vendor/github.com/golang/mock/mockgen
