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
			stack []func(*atomic.Uint32) run.Fn
		}
		results struct {
			work uint32
		}
	)

	ctx := t.Context()

	cases := map[string]struct {
		pl_testing.TestCase
		args    args
		results results
	}{
		"case00: success": {
			args: args{
				stack: []func(*atomic.Uint32) run.Fn{
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
				stack: []func(*atomic.Uint32) run.Fn{
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

			stack := x.SliceConvert(
				tc.args.stack,
				func(fnGetter func(*atomic.Uint32) run.Fn) run.Fn {
					if fnGetter == nil {
						return nil
					}

					return fnGetter(&counter)
				},
			)

			err := run.Group(ctx, stack...)

			require.Equal(t, tc.results.work, counter.Load())
			tc.CheckError(t, err)
		})
	}
}

func TestGroupSoft(t *testing.T) {
	pl_testing.Init(t)

	type (
		args struct {
			stack []func(*atomic.Uint32) run.Fn
		}
		results struct {
			work uint32
		}
	)

	ctx := t.Context()

	cases := map[string]struct {
		pl_testing.TestCase
		args    args
		results results
	}{
		"case00: success": {
			args: args{
				stack: []func(*atomic.Uint32) run.Fn{
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
				stack: []func(*atomic.Uint32) run.Fn{
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

			stack := x.SliceConvert(
				tc.args.stack,
				func(fnGetter func(*atomic.Uint32) run.Fn) run.Fn {
					if fnGetter == nil {
						return nil
					}

					return fnGetter(&counter)
				},
			)

			err := run.GroupSoft(ctx, stack...)

			require.Equal(t, tc.results.work, counter.Load())
			tc.CheckError(t, err)
		})
	}
}

func increment(c *atomic.Uint32) run.Fn {
	return func(context.Context) error {
		c.Add(1)

		return nil
	}
}

func fail(*atomic.Uint32) run.Fn {
	return func(context.Context) error {
		return io.EOF
	}
}
