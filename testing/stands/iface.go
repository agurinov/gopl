package stands

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type (
	State struct {
		Created bool
	}
	Interface interface {
		Up(t *testing.T) (created bool)
		Name() string
	}
)

func Init(
	t *testing.T,
	si ...Interface,
) map[string]State {
	t.Helper()

	states := make(map[string]State, len(si))

	for _, stand := range si {
		require.NotNil(t, stand)

		var (
			standName    = stand.Name()
			standCreated = stand.Up(t)
		)

		states[standName] = State{
			Created: standCreated,
		}
	}

	return states
}
