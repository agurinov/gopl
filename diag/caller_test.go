package diag_test

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/agurinov/gopl/diag"
	pl_testing "github.com/agurinov/gopl/testing"
)

type a struct{}

func (a) RunSimple() {}

func (a) Run(context.Context) error {
	return nil
}

func (a) Foo() string {
	return diag.CallerName(0)
}

func Foo() string {
	return diag.CallerName(0)
}

func Run(context.Context) error {
	return io.EOF
}

func TestCallerName(t *testing.T) {
	pl_testing.Init(t)

	type (
		args struct {
			f func() string
		}
		results struct {
			caller string
		}
	)

	cases := map[string]struct {
		args    args
		results results
		pl_testing.TestCase
	}{
		"case00: lambda": {
			args: args{
				f: func() string { return diag.CallerName(0) },
			},
			results: results{
				caller: "diag_test.TestCallerName.func1",
			},
		},
		"case01: func": {
			args: args{
				f: Foo,
			},
			results: results{
				caller: "diag_test.Foo",
			},
		},
		"case02: method": {
			args: args{
				f: a{}.Foo,
			},
			results: results{
				caller: "diag_test.a.Foo",
			},
		},
	}

	for name := range cases {
		tc := cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			caller := tc.args.f()
			require.Equal(t, tc.results.caller, caller)
		})
	}
}

func TestFunctionName(t *testing.T) {
	pl_testing.Init(t)

	type (
		args struct {
			f any
		}
		results struct {
			functionName string
		}
	)

	cases := map[string]struct {
		args    args
		results results
		pl_testing.TestCase
	}{
		"case00: lambda": {
			args: args{
				f: func(context.Context) error { return io.EOF },
			},
			results: results{
				functionName: "diag_test.TestFunctionName.func1",
			},
		},
		"case01: func": {
			args: args{
				f: Run,
			},
			results: results{
				functionName: "diag_test.Run",
			},
		},
		"case02: method": {
			args: args{
				f: a{}.Run,
			},
			results: results{
				functionName: "diag_test.a.Run",
			},
		},
		"case03: nil": {
			args: args{
				f: nil,
			},
			results: results{
				functionName: "nil",
			},
		},
		"case04: invalid": {
			args: args{
				f: 5,
			},
			results: results{
				functionName: "invalid",
			},
		},
	}

	for name := range cases {
		tc := cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			functionName := diag.FunctionName(tc.args.f)
			require.Equal(t, tc.results.functionName, functionName)
		})
	}
}
