package diag_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/agurinov/gopl/diag"
	pl_testing "github.com/agurinov/gopl/testing"
)

type a struct{}

func (a) Foo() string {
	return diag.CallerName()
}

func Foo() string {
	return diag.CallerName()
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
				f: diag.CallerName,
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
		name, tc := name, cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			caller := tc.args.f()
			require.Equal(t, tc.results.caller, caller)
		})
	}
}
