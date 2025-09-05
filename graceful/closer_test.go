package graceful_test

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/agurinov/gopl/graceful"
	pl_testing "github.com/agurinov/gopl/testing"
)

func simpleCloser()                          {}
func successCloser() error                   { return nil }
func errCloser() error                       { return io.EOF }
func successCtxCloser(context.Context) error { return nil }
func errCtxCloser(context.Context) error     { return io.ErrNoProgress }

func TestCloser_WaitForShutdown(t *testing.T) {
	pl_testing.Init(t)

	type (
		args struct {
			closers       []func()
			errClosers    []func() error
			ctxErrClosers []func(context.Context) error
		}
		results struct {
			errsIs []error
		}
	)

	cases := map[string]struct {
		args    args
		results results
		pl_testing.TestCase
	}{
		"case00: empty closer": {
			args:    args{},
			results: results{},
		},
		"case01: success closer": {
			args: args{
				closers:       []func(){simpleCloser},
				errClosers:    []func() error{successCloser},
				ctxErrClosers: []func(context.Context) error{successCtxCloser},
			},
			results: results{},
		},
		"case02: combined closer": {
			args: args{
				closers:       []func(){simpleCloser, simpleCloser},
				errClosers:    []func() error{successCloser, errCloser},
				ctxErrClosers: []func(context.Context) error{successCtxCloser, errCtxCloser},
			},
			results: results{
				errsIs: []error{io.EOF, io.ErrNoProgress},
			},
		},
	}

	for name := range cases {
		tc := cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			closer, err := graceful.NewCloser(
				graceful.WithLogger(zaptest.NewLogger(t)),
				graceful.WithTimeout(time.Second),
			)
			require.NoError(t, err)
			require.NotNil(t, closer)

			for _, fn := range tc.args.closers {
				closer.AddCloser(fn)
			}

			for _, fn := range tc.args.errClosers {
				closer.AddErrorCloser(fn)
			}

			for _, fn := range tc.args.ctxErrClosers {
				closer.AddContextErrorCloser(fn)
			}

			ctx, cancel := context.WithCancel(context.TODO())
			cancel()

			joinedErr := closer.WaitForShutdown(ctx)

			switch {
			case len(tc.results.errsIs) == 0:
				require.NoError(t, joinedErr)
			default:
				for _, targetErr := range tc.results.errsIs {
					require.ErrorIs(t, joinedErr, targetErr)
				}
			}
		})
	}
}
