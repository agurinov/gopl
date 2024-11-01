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
	jsonCfg struct {
		Logger    config.Logger
		Probes    config.Probes
		Graceful  config.Graceful
		GRPC      config.GRPC
		DebugHTTP config.HTTP       `json:"debug_http" yaml:"debug_http" `
		MergedMap map[string]string `json:"merged_map" yaml:"merged_map" `
	}
)

//go:embed testdata/config2.json
var sourceBytesJSON []byte

func TestParse_JSON(t *testing.T) {
	pl_testing.Init(t)

	type (
		args struct {
			sources []config.Source
		}
		results struct {
			cfg jsonCfg
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
					config.FromFile("testdata/config1.json"),
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
					config.FromFile("testdata/foobar.json"),
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
					config.FromFile("testdata/foobar.json", config.Silent),
					config.FromFile("testdata/config1.json"),
					config.FromFile("testdata/lolkek.json", config.Silent),
					config.FromBytes(sourceBytesJSON),
				},
			},
			results: results{
				cfg: jsonCfg{
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
		name, tc := name, cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			cfg, err := config.Parse[jsonCfg](ctx,
				config.JSON,
				tc.args.sources...,
			)
			tc.CheckError(t, err)
			require.Equal(t, tc.results.cfg, cfg)
		})
	}
}
