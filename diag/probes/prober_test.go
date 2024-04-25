package probes_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"

	"github.com/agurinov/gopl/diag/log"
	"github.com/agurinov/gopl/diag/probes"
	pl_testing "github.com/agurinov/gopl/testing"
)

func TestProber_New(t *testing.T) {
	type (
		args struct {
			opts []probes.Option
		}
	)

	cases := map[string]struct {
		args args
		pl_testing.TestCase
	}{
		"case00: empty options": {
			args: args{},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailAsErr: new(validator.ValidationErrors),
			},
		},
		"case01: all invalid": {
			args: args{
				opts: []probes.Option{
					probes.WithLogger(nil),
					probes.WithCheckInterval(time.Millisecond),
					probes.WithCheckTimeout(time.Millisecond),
					probes.WithProbe(probes.ProbeTypeReadiness, nil),
				},
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailAsErr: new(validator.ValidationErrors),
			},
		},
		"case03: all valid": {
			args: args{
				opts: []probes.Option{
					probes.WithLogger(log.NewZapTest(t)),
					probes.WithCheckInterval(time.Second),
					probes.WithCheckTimeout(time.Second),
					probes.WithProbe(
						probes.ProbeTypeReadiness,
						func(context.Context) error { return nil },
					),
				},
			},
		},
	}

	for name := range cases {
		name, tc := name, cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			prober, err := probes.New(tc.args.opts...)
			tc.CheckError(t, err)
			require.NotNil(t, prober)
		})
	}
}
