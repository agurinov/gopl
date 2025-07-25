package x_test

import (
	"strconv"
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
		tc := cases[name]

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

func TestMapConvert(t *testing.T) {
	pl_testing.Init(t)

	type (
		args struct {
			in   map[string]int
			mapF func(string, int) (int, string)
		}
		results struct {
			out map[int]string
		}
	)

	cases := map[string]struct {
		args    args
		results results
		pl_testing.TestCase
	}{
		"case00: typical": {
			args: args{
				in: map[string]int{
					"a":   10,
					"bb":  20,
					"ccc": 30,
				},
				mapF: func(k string, v int) (int, string) {
					return len(k), strconv.Itoa(v)
				},
			},
			results: results{
				out: map[int]string{
					1: "10",
					2: "20",
					3: "30",
				},
			},
		},
	}

	for name := range cases {
		tc := cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			out := x.MapConvert(
				tc.args.in,
				tc.args.mapF,
			)

			require.Equal(t, tc.results.out, out)
		})
	}
}

func TestMapFilter(t *testing.T) {
	pl_testing.Init(t)

	type (
		args struct {
			in   map[string]int
			useF func(string, int) bool
		}
		results struct {
			out map[string]int
		}
	)

	cases := map[string]struct {
		args    args
		results results
		pl_testing.TestCase
	}{
		"case00: filter values greater than 1": {
			args: args{
				in: map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				useF: func(_ string, v int) bool { return v > 1 },
			},
			results: results{
				out: map[string]int{
					"bar": 2,
					"baz": 3,
				},
			},
		},
		"case01: filter by key length": {
			args: args{
				in: map[string]int{
					"a":    1,
					"ab":   2,
					"abc":  3,
					"abcd": 4,
				},
				useF: func(k string, _ int) bool { return len(k) > 2 },
			},
			results: results{
				out: map[string]int{
					"abc":  3,
					"abcd": 4,
				},
			},
		},
		"case02: no elements satisfy predicate": {
			args: args{
				in: map[string]int{
					"x": 0,
					"y": 0,
				},
				useF: func(_ string, v int) bool { return v > 1 },
			},
			results: results{
				out: map[string]int{},
			},
		},
	}

	for name := range cases {
		tc := cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			out := x.MapFilter(
				tc.args.in,
				tc.args.useF,
			)

			require.Equal(t, tc.results.out, out)
		})
	}
}
