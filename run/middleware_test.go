package run_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/agurinov/gopl/run"
	pl_testing "github.com/agurinov/gopl/testing"
)

func TestMiddlewares_Handler(t *testing.T) {
	pl_testing.Init(t)

	type Handler = func([]string) []string

	type (
		args struct {
			mws run.Middlewares[Handler]
		}
		results struct {
			chain []string
		}
	)

	var (
		root = func(in []string) []string {
			return append(in, "root")
		}
		wrap = func(s string) run.Middleware[Handler] {
			return func(h Handler) Handler {
				return func(in []string) []string {
					in = append(in, s+"_before")
					in = h(in)
					in = append(in, s+"_after")

					return in
				}
			}
		}
	)

	cases := map[string]struct {
		args    args
		results results
		pl_testing.TestCase
	}{
		"case00: success": {
			args: args{
				mws: run.Middlewares[Handler]{
					wrap("outer"),
					wrap("inner"),
				},
			},
			results: results{
				chain: []string{
					"outer_before",
					"inner_before",
					"root",
					"inner_after",
					"outer_after",
				},
			},
		},
	}

	for name := range cases {
		tc := cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			handler := tc.args.mws.Handler(root)
			require.NotNil(t, handler)

			chainSize := 2*len(tc.args.mws) + 1

			chain := handler(make([]string, 0, chainSize))

			require.Equal(t, tc.results.chain, chain)
		})
	}
}
