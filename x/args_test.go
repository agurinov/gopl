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
		pl_testing.TestCase
		results results
		args    args
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

func TestSafeAssert(t *testing.T) {
	pl_testing.Init(t)

	// require.Equal(t, int64(5), x.SafeAssert[int64](uint(5)))

	require.Equal(t, "", x.SafeAssert[string](5))
	require.Equal(t, "", x.SafeAssert[string](true))
	require.Equal(t, "", x.SafeAssert[string](nil))
	require.Equal(t, "foobar", x.SafeAssert[string]("foobar"))
}

func TestParseNumber(t *testing.T) {
	pl_testing.Init(t)

	outInt, err := x.ParseNumber[int]("   ")
	require.NoError(t, err)
	require.Equal(t, int(0), outInt)

	outInt32, err := x.ParseNumber[int32]("100,500")
	require.NoError(t, err)
	require.Equal(t, int32(100500), outInt32)

	outInt64, err := x.ParseNumber[int64]("100,500.50")
	require.NoError(t, err)
	require.Equal(t, int64(100500), outInt64)

	outFloat64, err := x.ParseNumber[float64]("100,500.50")
	require.NoError(t, err)
	require.InDelta(t, float64(100500.50), outFloat64, 0)

	NaN, err := x.ParseNumber[float32]("lolkek")
	require.Error(t, err)
	require.InDelta(t, float32(0), NaN, 0)
}
