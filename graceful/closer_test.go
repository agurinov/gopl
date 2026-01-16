package graceful_test

import (
	"context"
	"errors"
	"io"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/agurinov/gopl/graceful"
	"github.com/agurinov/gopl/run"
	pl_testing "github.com/agurinov/gopl/testing"
)

type db struct {
	closed atomic.Bool
}

type svc struct {
	db *db
}

func simpleCloser()                          {}
func successCloser() error                   { return nil }
func errCloser() error                       { return io.EOF }
func successCtxCloser(context.Context) error { return nil }
func errCtxCloser(context.Context) error     { return io.ErrNoProgress }

func (d *db) CommitAll() error {
	if d.closed.Load() {
		return errors.New("db is closed")
	}

	return nil
}

func (d *db) Close() {
	d.closed.Store(true)
}

func (s svc) Close() error {
	return s.db.CommitAll()
}

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
				graceful.WithCloserLogger(zaptest.NewLogger(t)),
				graceful.WithCloserTimeout(time.Second),
			)
			require.NoError(t, err)
			require.NotNil(t, closer)

			for _, fn := range tc.args.closers {
				closer.AddCloser(run.SimpleFn(fn))
			}

			for _, fn := range tc.args.errClosers {
				closer.AddCloser(run.ErrorFn(fn))
			}

			for _, fn := range tc.args.ctxErrClosers {
				closer.AddCloser(fn)
			}

			ctx, cancel := context.WithCancel(t.Context())
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

func TestCloser_Waves(t *testing.T) {
	pl_testing.Init(t)

	closer, err := graceful.NewCloser(
		graceful.WithCloserLogger(zaptest.NewLogger(t)),
		graceful.WithCloserTimeout(time.Second),
	)
	require.NoError(t, err)
	require.NotNil(t, closer)

	database := new(db)
	closer.AddCloser(run.SimpleFn(database.Close))

	service := svc{db: database}
	closer.AddCloser(
		run.ErrorFn(service.Close),
		graceful.InFirstWave(),
	)

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	joinedErr := closer.WaitForShutdown(ctx)
	require.NoError(t, joinedErr)
}
