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
		isDebug       = tc.Debug
		isIntegration = len(si) > 0
	)

	var (
		needDebug       = envvars.GDebug.Present()
		needIntegration = !testing.Short()
		needParallel    = !needDebug
	)

	switch {
	case tc.Skip:
		t.Skip("tc skipped: explicit skip flag")
	case tc.Fail:
		t.Fail()
	case !needIntegration && isIntegration:
		t.Skip("tc skipped: integration test during unit mode")
	case tc.root:
	case needDebug && !isDebug:
		t.Skip("tc skipped: not debuggable during " + envvars.GDebug.String())
	case needIntegration && !isIntegration:
		t.Skip("tc skipped: unit test during integration mode")
	}

	if needParallel {
		t.Parallel()
	}

	return stands.Init(t, si...)
}
