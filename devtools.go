//go:build neverbuild

package main

// TODO: Improve this pattern
// https://www.jvt.me/posts/2025/01/27/go-tools-124/

import (
	_ "github.com/vektra/mockery/v3"
	_ "golang.org/x/tools/cmd/stringer"
)
