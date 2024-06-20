//go:build py38

// https://docs.python.org/3.8/c-api/index.html
package cpy

// #cgo pkg-config: python-3.8-embed
import "C"

const (
	MinorVersion = "3.8"
)
