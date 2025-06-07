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
		MyInstance struct{}
	)

	myErr := ddd.NotFoundError[MyInstance]{}

	cases := map[string]struct {
		args    args
		results results
		pl_testing.TestCase
	}{
		"case00: is not found (generic)": {
			args: args{
				err: myErr,
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
	}

	for name := range cases {
		name, tc := name, cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			require.Equal(t, tc.results.isStd, errors.Is(tc.args.err, myErr))
			require.Equal(t, tc.results.isNotFound, ddd.IsNotFound(tc.args.err))
		})
	}
}
