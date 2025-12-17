package graceful_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/agurinov/gopl/graceful"
	pl_testing "github.com/agurinov/gopl/testing"
)

type a struct{}

func (a) RunSimple()                {}
func (a) RunError() error           { return nil }
func (a) Run(context.Context) error { return nil }

func TestClosure_String(t *testing.T) {
	pl_testing.Init(t)

	type (
		args struct {
			f graceful.Closure
		}
		results struct {
			asString string
		}
	)

	cases := map[string]struct {
		args    args
		results results
		pl_testing.TestCase
	}{
		"case00: simple closure": {
			args: args{
				f: graceful.SimpleClosure(a{}.RunSimple),
			},
			results: results{
				asString: "graceful_test.a.RunSimple",
			},
			TestCase: pl_testing.TestCase{
				Skip: true,
			},
		},
		"case01: error closure": {
			args: args{
				f: graceful.ErrorClosure(a{}.RunError),
			},
			results: results{
				asString: "graceful_test.a.RunError",
			},
			TestCase: pl_testing.TestCase{
				Skip: true,
			},
		},
		"case02: closure": {
			args: args{
				f: a{}.Run,
			},
			results: results{
				asString: "graceful_test.a.Run",
			},
		},
	}

	for name := range cases {
		tc := cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			asString := tc.args.f.String()
			require.Equal(t, tc.results.asString, asString)
		})
	}
}
