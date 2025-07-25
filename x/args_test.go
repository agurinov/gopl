package x_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	pl_testing "github.com/agurinov/gopl/testing"
	"github.com/agurinov/gopl/x"
)

func TestEmptyIf(t *testing.T) {
	pl_testing.Init(t)

	type (
		args struct {
			in    string
			empty []string
		}
		results struct {
			out string
		}
	)

	cases := map[string]struct {
		args    args
		results results
		pl_testing.TestCase
	}{
		"case00: no empty variants": {
			args: args{
				in: "foobar",
			},
			results: results{
				out: "foobar",
			},
		},
		"case01: json": {
			args: args{
				in: "null",
				empty: []string{
					"[]",
					"{}",
					"null",
				},
			},
			results: results{
				out: "",
			},
		},
	}

	for name := range cases {
		tc := cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			out := x.EmptyIf(tc.args.in, tc.args.empty...)
			require.Equal(t, tc.results.out, out)
		})
	}
}
