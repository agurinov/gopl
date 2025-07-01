package x_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	pl_testing "github.com/agurinov/gopl/testing"
	"github.com/agurinov/gopl/x"
)

func TestUnique(t *testing.T) {
	pl_testing.Init(t)

	type (
		args struct {
			in []string
		}
		results struct {
			unique []string
		}
	)

	cases := map[string]struct {
		pl_testing.TestCase
		args    args
		results results
	}{
		"case00: without duplicates": {
			args: args{
				in: []string{"a", "b", "c"},
			},
			results: results{
				unique: []string{"a", "b", "c"},
			},
		},
		"case01: with duplicates": {
			args: args{
				in: []string{"a", "b", "c", "c", "a"},
			},
			results: results{
				unique: []string{"a", "b", "c"},
			},
		},
	}

	for name := range cases {
		name, tc := name, cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			unique := x.Unique(tc.args.in)
			require.Equal(t, tc.results.unique, unique)
		})
	}
}

func TestFilterOutEmpty(t *testing.T) {
	pl_testing.Init(t)

	type (
		args struct {
			in []string
		}
		results struct {
			filtered []string
		}
	)

	cases := map[string]struct {
		pl_testing.TestCase
		args    args
		results results
	}{
		"case00: without empty": {
			args: args{
				in: []string{"a", "b", "c"},
			},
			results: results{
				filtered: []string{"a", "b", "c"},
			},
		},
		"case01: all empty": {
			args: args{
				in: []string{"", "", ""},
			},
			results: results{
				filtered: nil,
			},
		},
		"case02: mixed": {
			args: args{
				in: []string{"", "a", "b"},
			},
			results: results{
				filtered: []string{"a", "b"},
			},
		},
	}

	for name := range cases {
		name, tc := name, cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			filtered := x.FilterOutEmpty(tc.args.in)
			require.Equal(t, tc.results.filtered, filtered)
		})
	}
}

func TestLast(t *testing.T) {
	pl_testing.Init(t)

	type (
		args struct {
			in []string
		}
		results struct {
			last string
		}
	)

	cases := map[string]struct {
		pl_testing.TestCase
		args    args
		results results
	}{
		"case00: nil": {
			args: args{
				in: nil,
			},
			results: results{
				last: "",
			},
		},
		"case01: empty": {
			args: args{
				in: []string{},
			},
			results: results{
				last: "",
			},
		},
		"case02: mixed": {
			args: args{
				in: []string{"a", "b", "c"},
			},
			results: results{
				last: "c",
			},
		},
	}

	for name := range cases {
		name, tc := name, cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			last := x.Last(tc.args.in)
			require.Equal(t, tc.results.last, last)
		})
	}
}

func TestCoalesce(t *testing.T) {
	pl_testing.Init(t)

	type (
		args struct {
			in []string
		}
		results struct {
			coalesce string
		}
	)

	cases := map[string]struct {
		pl_testing.TestCase
		args    args
		results results
	}{
		"case00: nil": {
			args: args{
				in: nil,
			},
			results: results{
				coalesce: "",
			},
		},
		"case01: empty": {
			args: args{
				in: []string{},
			},
			results: results{
				coalesce: "",
			},
		},
		"case02: mixed": {
			args: args{
				in: []string{"", "a", "b"},
			},
			results: results{
				coalesce: "a",
			},
		},
	}

	for name := range cases {
		name, tc := name, cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			coalesce := x.Coalesce(tc.args.in...)
			require.Equal(t, tc.results.coalesce, coalesce)
		})
	}
}

func TestSliceToMap(t *testing.T) {
	pl_testing.Init(t)

	type (
		s struct {
			UUID uuid.UUID
			S    string
		}
		args struct {
			in []s
		}
		results struct {
			out map[uuid.UUID]string
		}
	)

	cases := map[string]struct {
		pl_testing.TestCase
		args    args
		results results
	}{
		"case00: nil": {
			args: args{
				in: []s{
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
			results: results{
				out: map[uuid.UUID]string{
					uuid.Must(uuid.Parse("9B3C6872-F8C7-46D6-8376-DF39D61EB099")): "foo",
					uuid.Must(uuid.Parse("5408546F-FEEB-4348-A2B4-AB73C0D8E775")): "bar",
				},
			},
		},
	}

	for name := range cases {
		name, tc := name, cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			out := x.SliceToMap(
				tc.args.in,
				func(in s) (uuid.UUID, string) {
					return in.UUID, in.S
				},
			)
			require.EqualValues(t, tc.results.out, out)
		})
	}
}

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

func TestPaginate(t *testing.T) {
	pl_testing.Init(t)

	type (
		args struct {
			in     []string
			limit  uint
			offset uint
		}
		results struct {
			out []string
		}
	)

	cases := map[string]struct {
		pl_testing.TestCase
		args    args
		results results
	}{
		"case00: zero limit": {
			args: args{
				in:    []string{"a", "b", "c"},
				limit: 0,
			},
			results: results{
				out: nil,
			},
		},
		"case01: limit over len": {
			args: args{
				in:    []string{"a", "b", "c"},
				limit: 100,
			},
			results: results{
				out: []string{"a", "b", "c"},
			},
		},
		"case02: limit less len": {
			args: args{
				in:    []string{"a", "b", "c"},
				limit: 2,
			},
			results: results{
				out: []string{"a", "b"},
			},
		},
		"case03: offset overflow": {
			args: args{
				in:     []string{"a", "b", "c"},
				offset: 10,
				limit:  10,
			},
			results: results{
				out: nil,
			},
		},
		"case04: mixed": {
			args: args{
				in:     []string{"a", "b", "c"},
				offset: 1,
				limit:  1,
			},
			results: results{
				out: []string{"b"},
			},
		},
	}

	for name := range cases {
		name, tc := name, cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			out := x.Paginate(
				tc.args.in,
				tc.args.limit,
				tc.args.offset,
			)
			require.Equal(t, tc.results.out, out)
		})
	}
}

func TestFlatten(t *testing.T) {
	pl_testing.Init(t)

	type (
		args struct {
			in [][]string
		}
		results struct {
			out []string
		}
	)

	cases := map[string]struct {
		pl_testing.TestCase
		args    args
		results results
	}{
		"case00: mixed": {
			args: args{
				in: [][]string{
					{"a", "b"},
					{"c", "d"},
				},
			},
			results: results{
				out: []string{"a", "b", "c", "d"},
			},
		},
		"case01: empty": {
			args: args{
				in: [][]string{},
			},
			results: results{
				out: []string{},
			},
		},
		"case02: nil": {
			args: args{
				in: nil,
			},
			results: results{
				out: nil,
			},
		},
	}

	for name := range cases {
		name, tc := name, cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			out := x.Flatten(tc.args.in)
			require.Equal(t, tc.results.out, out)
		})
	}
}
