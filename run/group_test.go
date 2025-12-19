package run_test

import (
	"context"
	"io"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/agurinov/gopl/run"
	pl_testing "github.com/agurinov/gopl/testing"
	"github.com/agurinov/gopl/x"
)

func TestGroup(t *testing.T) {
	pl_testing.Init(t)

	type (
		args struct {
			stack []func(*atomic.Uint32) run.Closure
		}
		results struct {
			work uint32
		}
	)

	ctx := context.TODO()

	cases := map[string]struct {
		args    args
		results results
		pl_testing.TestCase
	}{
		"case00: success": {
			args: args{
				stack: []func(*atomic.Uint32) run.Closure{
					nil, increment,
					nil, increment,
					nil, increment,
				},
			},
			results: results{
				work: 3,
			},
		},
		"case01: with error": {
			args: args{
				stack: []func(*atomic.Uint32) run.Closure{
					nil, increment,
					nil, fail,
					nil, increment,
					nil, fail,
					nil, increment,
					nil, fail,
				},
			},
			results: results{
				work: 3,
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: io.EOF,
			},
		},
	}

	for name := range cases {
		tc := cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			var counter atomic.Uint32

			closures := x.SliceConvert(
				tc.args.stack,
				func(closureGetter func(*atomic.Uint32) run.Closure) run.Closure {
					if closureGetter == nil {
						return nil
					}

					return closureGetter(&counter)
				},
			)

			err := run.Group(ctx, closures...)

			require.Equal(t, tc.results.work, counter.Load())
			tc.CheckError(t, err)
		})
	}
}

func TestGroupSoft(t *testing.T) {
	pl_testing.Init(t)

	type (
		args struct {
			stack []func(*atomic.Uint32) run.Closure
		}
		results struct {
			work uint32
		}
	)

	ctx := context.TODO()

	cases := map[string]struct {
		args    args
		results results
		pl_testing.TestCase
	}{
		"case00: success": {
			args: args{
				stack: []func(*atomic.Uint32) run.Closure{
					nil, increment,
					nil, increment,
				},
			},
			results: results{
				work: 2,
			},
		},
		"case01: with error": {
			args: args{
				stack: []func(*atomic.Uint32) run.Closure{
					nil, increment,
					nil, fail,
					nil, increment,
					nil, fail,
					nil, increment,
					nil, fail,
				},
			},
			results: results{
				work: 3,
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: io.EOF,
			},
		},
	}

	for name := range cases {
		tc := cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			var counter atomic.Uint32

			closures := x.SliceConvert(
				tc.args.stack,
				func(closureGetter func(*atomic.Uint32) run.Closure) run.Closure {
					if closureGetter == nil {
						return nil
					}

					return closureGetter(&counter)
				},
			)

			err := run.GroupSoft(ctx, closures...)

			require.Equal(t, tc.results.work, counter.Load())
			tc.CheckError(t, err)
		})
	}
}

func increment(c *atomic.Uint32) run.Closure {
	return func(context.Context) error {
		c.Add(1)

		return nil
	}
}

func fail(*atomic.Uint32) run.Closure {
	return func(context.Context) error {
		return io.EOF
	}
}
