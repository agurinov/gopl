package x_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	pl_testing "github.com/agurinov/gopl/testing"
	"github.com/agurinov/gopl/x"
)

func TestDecartDuo(t *testing.T) {
	pl_testing.Init(t)

	type (
		A struct {
			s string
			i int
		}
		B struct {
			b bool
		}
		E struct {
			a A
			B
		}
	)

	type (
		args struct {
			as []A
			bs []B
		}
		results struct {
			e []E
		}
	)

	cases := map[string]struct {
		args    args
		results results
		pl_testing.TestCase
	}{
		"case00: empty": {
			args: args{},
			results: results{
				e: []E{},
			},
		},
		"case01: only first dimension": {
			args: args{
				as: []A{
					{s: "s1", i: 1},
					{s: "s2", i: 2},
					{s: "s3", i: 3},
				},
			},
			results: results{
				e: []E{},
			},
		},
		"case02: only second dimension": {
			args: args{
				bs: []B{
					{b: true},
					{b: false},
				},
			},
			results: results{
				e: []E{},
			},
		},
		"case03: two dimensions": {
			args: args{
				as: []A{
					{s: "s1", i: 1},
					{s: "s2", i: 2},
					{s: "s3", i: 3},
				},
				bs: []B{
					{b: true},
					{b: false},
				},
			},
			results: results{
				e: []E{
					{
						a: A{s: "s1", i: 1},
						B: B{b: true},
					},
					{
						a: A{s: "s1", i: 1},
						B: B{b: false},
					},
					{
						a: A{s: "s2", i: 2},
						B: B{b: true},
					},
					{
						a: A{s: "s2", i: 2},
						B: B{b: false},
					},
					{
						a: A{s: "s3", i: 3},
						B: B{b: true},
					},
					{
						a: A{s: "s3", i: 3},
						B: B{b: false},
					},
				},
			},
		},
	}

	for name := range cases {
		tc := cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			e := x.DecartDuo(
				tc.args.as,
				tc.args.bs,
				func(a A, b B) E {
					return E{
						a: a,
						B: b,
					}
				},
			)
			require.Equal(t, tc.results.e, e)
		})
	}
}
