package config_test

import (
	"context"
	_ "embed"
	"os"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"

	"github.com/agurinov/gopl/config"
	pl_testing "github.com/agurinov/gopl/testing"
)

type (
	yamlCfg struct {
		Logger    config.Logger
		Probes    config.Probes
		Graceful  config.Graceful
		GRPC      config.GRPC
		DebugHTTP config.HTTP       `json:"debug_http" yaml:"debug_http" `
		MergedMap map[string]string `json:"merged_map" yaml:"merged_map" `
	}
)

//go:embed testdata/config2.yaml
var sourceBytesYAML []byte

func TestParse_YAML(t *testing.T) {
	pl_testing.Init(t)

	type (
		args struct {
			sources []config.Source
		}
		results struct {
			cfg yamlCfg
		}
	)

	ctx := context.TODO()

	cases := map[string]struct {
		args    args
		results results
		pl_testing.TestCase
	}{
		"case00: no sources": {
			args: args{
				sources: []config.Source{nil, nil, nil},
			},
			results: results{},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailAsErr: new(validator.ValidationErrors),
			},
		},
		"case01: file": {
			args: args{
				sources: []config.Source{
					config.FromFile("testdata/config1.yaml"),
				},
			},
			results: results{},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailAsErr: new(validator.ValidationErrors),
			},
		},
		"case02: file doesnt exists strict": {
			args: args{
				sources: []config.Source{
					config.FromFile("testdata/foobar.yaml"),
				},
			},
			results: results{},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: os.ErrNotExist,
			},
		},
		"case03: file + bytes = ok": {
			args: args{
				sources: []config.Source{
					config.FromFile("testdata/foobar.yaml", config.Silent),
					config.FromFile("testdata/config1.yaml"),
					config.FromFile("testdata/lolkek.yaml", config.Silent),
					config.FromBytes(sourceBytesYAML),
				},
			},
			results: results{
				cfg: yamlCfg{
					Logger: config.Logger{
						Format: "json",
						Level:  "info",
					},
					Probes: config.Probes{
						CheckInterval: 3 * time.Second,
						CheckTimeout:  5 * time.Second,
					},
					Graceful: config.Graceful{
						ShutdownTimeout: 2 * time.Second,
						CloseTimeout:    time.Second,
					},
					GRPC: config.GRPC{
						Port:             9090,
						MaxRequestBytes:  100,
						MaxResponseBytes: 100,
					},
					DebugHTTP: config.HTTP{
						Port: 8081,
					},
					MergedMap: map[string]string{
						"foo": "bar2",
						"bar": "baz",
						"lol": "kek",
					},
				},
			},
		},
	}

	for name := range cases {
		tc := cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			cfg, err := config.Parse[yamlCfg](ctx,
				config.YAML,
				tc.args.sources...,
			)
			tc.CheckError(t, err)
			require.Equal(t, tc.results.cfg, cfg)
		})
	}
}
