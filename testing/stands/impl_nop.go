package stands

import (
	"testing"
)

type nop struct{}

var Nop = nop{}

func (nop) Name() string       { return "nop" }
func (nop) Up(*testing.T) bool { return true }
