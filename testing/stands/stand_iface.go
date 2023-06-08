package stands

import "testing"

type (
	State struct {
		Created bool
	}
	Interface interface {
		Up(t *testing.T) (created bool)
		Name() string
	}
)
