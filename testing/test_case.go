package testing

import (
	"testing"

	"github.com/agurinov/gopl/env/envvars"
	"github.com/agurinov/gopl/testing/stands"
)

type TestCase struct {
	MustFailIsErr error
	MustFailAsErr any
	MustFail      bool

	root  bool
	Skip  bool
	Debug bool
	Fail  bool
}

func Init(t *testing.T, si ...stands.Interface) map[string]stands.State {
	t.Helper()

	tc := TestCase{
		root: true,
	}

	return tc.Init(t, si...)
}

func (tc TestCase) Init(t *testing.T, si ...stands.Interface) map[string]stands.State {
	t.Helper()

	var (
		isIntegration = len(si) > 0
		isDebug       = tc.Debug || tc.root
		needDebug     = envvars.GDebug.Present()
		needParallel  = !needDebug
	)

	switch {
	case tc.Skip:
		t.Skip("tc skipped: explicit skip flag")
	case needDebug && !isDebug:
		t.Skip("tc skipped: not debuggable during " + envvars.GDebug.String())
	case testing.Short() && isIntegration:
		t.Skip("tc skipped: integration test during -testing.short mode")
	case tc.Fail:
		t.Fail()
	}

	if needParallel {
		t.Parallel()
	}

	states := make(map[string]stands.State, len(si))

	for _, stand := range si {
		var (
			standName    = stand.Name()
			standCreated = stand.Up(t)
		)

		states[standName] = stands.State{
			Created: standCreated,
		}
	}

	return states
}
