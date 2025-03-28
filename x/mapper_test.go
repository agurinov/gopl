package x_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	pl_testing "github.com/agurinov/gopl/testing"
	"github.com/agurinov/gopl/x"
)

func TestSliceMapError_UUID(t *testing.T) {
	pl_testing.Init(t)

	type (
		args struct {
			in []string
		}
		results struct {
			out []uuid.UUID
		}
	)

	cases := map[string]struct {
		args    args
		results results
		pl_testing.TestCase
	}{
		"case00: nil": {
			args: args{
				in: nil,
			},
			results: results{
				out: nil,
			},
		},
		"case01: empty": {
			args: args{
				in: []string{},
			},
			results: results{
				out: nil,
			},
		},
		"case02: some invalid uuid": {
			args: args{
				in: []string{
					"2dbd6c1c-106a-495c-accb-b900c232c8f2",
					"foo",
					"ea7d34b3-6602-4fe0-b680-3089f7dee046",
				},
			},
			results: results{
				out: nil,
			},
			TestCase: pl_testing.TestCase{
				MustFail: true,
			},
		},
		"case03: all valid": {
			args: args{
				in: []string{
					"2dbd6c1c-106a-495c-accb-b900c232c8f2",
					"ea7d34b3-6602-4fe0-b680-3089f7dee046",
				},
			},
			results: results{
				out: []uuid.UUID{
					uuid.Must(uuid.Parse("2dbd6c1c-106a-495c-accb-b900c232c8f2")),
					uuid.Must(uuid.Parse("ea7d34b3-6602-4fe0-b680-3089f7dee046")),
				},
			},
		},
	}

	for name := range cases {
		name, tc := name, cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			out, err := x.SliceMapError(
				tc.args.in,
				uuid.Parse,
			)
			tc.CheckError(t, err)
			require.Equal(t, tc.results.out, out)
		})
	}
}
