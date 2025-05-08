package x_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	pl_testing "github.com/agurinov/gopl/testing"
	"github.com/agurinov/gopl/x"
)

func TestMapToSlice(t *testing.T) {
	pl_testing.Init(t)

	type (
		s struct {
			UUID uuid.UUID
			S    string
		}
		args struct {
			in map[uuid.UUID]string
		}
		results struct {
			out []s
		}
	)

	cases := map[string]struct {
		args    args
		results results
		pl_testing.TestCase
	}{
		"case00: mixed": {
			args: args{
				in: map[uuid.UUID]string{
					uuid.Must(uuid.Parse("9B3C6872-F8C7-46D6-8376-DF39D61EB099")): "foo",
					uuid.Must(uuid.Parse("5408546F-FEEB-4348-A2B4-AB73C0D8E775")): "bar",
				},
			},
			results: results{
				out: []s{
					{
						UUID: uuid.Must(uuid.Parse("9B3C6872-F8C7-46D6-8376-DF39D61EB099")),
						S:    "foo",
					},
					{
						UUID: uuid.Must(uuid.Parse("5408546F-FEEB-4348-A2B4-AB73C0D8E775")),
						S:    "bar",
					},
				},
			},
		},
	}

	for name := range cases {
		name, tc := name, cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			out := x.MapToSlice(
				tc.args.in,
				func(k uuid.UUID, v string) s {
					return s{
						UUID: k,
						S:    v,
					}
				},
			)
			require.ElementsMatch(t, tc.results.out, out)
		})
	}
}
