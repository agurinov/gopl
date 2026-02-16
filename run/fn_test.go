package run_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/agurinov/gopl/diag"
	"github.com/agurinov/gopl/run"
	pl_testing "github.com/agurinov/gopl/testing"
)

type a struct{}

func (a) RunSimple()                {}
func (a) RunError() error           { return nil }
func (a) Run(context.Context) error { return nil }

func TestFn_String(t *testing.T) {
	pl_testing.Init(t)

	type (
		args struct {
			f run.Fn
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
		"case00: simple fn": {
			args: args{
				f: run.SimpleFn(a{}.RunSimple),
			},
			results: results{
				asString: "run_test.a.RunSimple",
			},
			TestCase: pl_testing.TestCase{
				Skip: true,
			},
		},
		"case01: error fn": {
			args: args{
				f: run.ErrorFn(a{}.RunError),
			},
			results: results{
				asString: "run_test.a.RunError",
			},
			TestCase: pl_testing.TestCase{
				Skip: true,
			},
		},
		"case02: fn": {
			args: args{
				f: a{}.Run,
			},
			results: results{
				asString: "run_test.a.Run",
			},
		},
	}

	for name := range cases {
		tc := cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			asString := diag.FunctionName(tc.args.f)
			require.Equal(t, tc.results.asString, asString)
		})
	}
}
