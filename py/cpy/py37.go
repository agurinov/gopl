//go:build ignore

// https://docs.python.org/3.7/c-api/index.html
package cpy

// #cgo pkg-config: python-3.7
import "C"

const (
	MinorVersion = "3.7"
)
