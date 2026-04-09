package ddd_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/agurinov/gopl/ddd"
	pl_testing "github.com/agurinov/gopl/testing"
)

func TestIsNotFound(t *testing.T) {
	pl_testing.Init(t)

	type (
		args struct {
			err error
		}
		results struct {
			isStd      bool
			isNotFound bool
		}
		Foo struct{ Foo string }
		Bar struct{ Bar string }
	)

	var (
		fooErr ddd.NotFoundError[Foo]
		barErr ddd.NotFoundError[Bar]
	)

	cases := map[string]struct {
		args args
		pl_testing.TestCase
		results results
	}{
		"case00: generic is not found same constant": {
			args: args{
				err: fooErr,
			},
			results: results{
				isStd:      true,
				isNotFound: true,
			},
		},
		"case01: nil error": {
			args: args{
				err: nil,
			},
			results: results{
				isStd:      false,
				isNotFound: false,
			},
		},
		"case02: other error": {
			args: args{
				err: errors.New("other"),
			},
			results: results{
				isStd:      false,
				isNotFound: false,
			},
		},
		"case03: generic is not found same type": {
			args: args{
				err: ddd.NotFoundError[Foo]{},
			},
			results: results{
				isStd:      true,
				isNotFound: true,
			},
		},
		"case04: generic is not found another type": {
			args: args{
				err: barErr,
			},
			results: results{
				isStd:      false,
				isNotFound: true,
			},
		},
	}

	for name := range cases {
		tc := cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			require.Equal(t, tc.results.isStd, errors.Is(tc.args.err, fooErr))
			require.Equal(t, tc.results.isNotFound, ddd.IsNotFound(tc.args.err))
		})
	}
}
