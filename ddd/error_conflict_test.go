package ddd_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/agurinov/gopl/ddd"
	pl_testing "github.com/agurinov/gopl/testing"
)

func TestIsConflict(t *testing.T) {
	pl_testing.Init(t)

	type (
		args struct {
			err error
		}
		results struct {
			isStd      bool
			isConflict bool
		}
		Foo struct{ Foo string }
		Bar struct{ Bar string }
	)

	var (
		fooErr ddd.ConflictError[Foo]
		barErr ddd.ConflictError[Bar]
	)

	cases := map[string]struct {
		args    args
		results results
		pl_testing.TestCase
	}{
		"case00: generic conflict same constant": {
			args: args{
				err: fooErr,
			},
			results: results{
				isStd:      true,
				isConflict: true,
			},
		},
		"case01: nil error": {
			args: args{
				err: nil,
			},
			results: results{
				isStd:      false,
				isConflict: false,
			},
		},
		"case02: other error": {
			args: args{
				err: errors.New("other"),
			},
			results: results{
				isStd:      false,
				isConflict: false,
			},
		},
		"case03: generic conflict same type": {
			args: args{
				err: ddd.ConflictError[Foo]{},
			},
			results: results{
				isStd:      true,
				isConflict: true,
			},
		},
		"case04: generic conflict another type": {
			args: args{
				err: barErr,
			},
			results: results{
				isStd:      false,
				isConflict: true,
			},
		},
	}

	for name := range cases {
		tc := cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			require.Equal(t, tc.results.isStd, errors.Is(tc.args.err, fooErr))
			require.Equal(t, tc.results.isConflict, ddd.IsConflict(tc.args.err))
		})
	}
}
