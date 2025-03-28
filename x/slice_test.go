package x_test

import (
	"testing"

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
