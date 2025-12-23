//go:build neverbuild

package main

// TODO: Improve this pattern
// https://www.jvt.me/posts/2025/01/27/go-tools-124/

import (
	_ "go.uber.org/mock/mockgen"
	_ "golang.org/x/tools/cmd/stringer"
)
