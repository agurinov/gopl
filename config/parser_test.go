package config_test

import (
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
	myStruct struct {
		Key   string
		Value string
	}
	myCfg struct {
		MapScalar   map[string]string   `json:"map_scalar_j" yaml:"map_scalar_y"`
		MapStruct   map[string]myStruct `json:"map_struct_j" yaml:"map_struct_y"`
		Logger      config.Logger       `json:"logger_j" yaml:"logger_y"`
		DebugHTTP   config.HTTP         `json:"debug_http_j" yaml:"debug_http_y"`
		ArrayScalar []string            `json:"array_scalar_j" yaml:"array_scalar_y"`
		ArrayStruct []myStruct          `json:"array_struct_j" yaml:"array_struct_y"`
		GRPC        config.GRPC         `json:"grpc_j" yaml:"grpc_y"`
		Probes      config.Probes       `json:"probes_j" yaml:"probes_y"`
		Graceful    config.Graceful     `json:"graceful_j" yaml:"graceful_y"`
	}
)

var (
	//go:embed testdata/config2.yaml
	config2YAML []byte

	//go:embed testdata/config2.json
	config2JSON []byte
)

func TestParse(t *testing.T) {
	pl_testing.Init(t)

	type (
		args struct {
			sources []config.Source
			flags   config.Flags
		}
		results struct {
			cfg myCfg
		}
	)

	ctx := t.Context()

	t.Run("JSON", func(t *testing.T) {
		pl_testing.Init(t)

		cases := map[string]struct {
			pl_testing.TestCase
			args    args
			results results
		}{
			"case00: no sources": {
				args: args{
					flags:   config.JSON,
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
					flags: config.JSON,
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
					flags: config.JSON,
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
					flags: config.JSON,
					sources: []config.Source{
						config.FromFile("testdata/foobar.json", config.Silent),
						config.FromFile("testdata/config1.json"),
						config.FromFile("testdata/lolkek.json", config.Silent),
						config.FromBytes(config2JSON),
					},
				},
				results: results{
					cfg: myCfg{
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
						MapScalar: map[string]string{
							"foo": "bar2",
							"bar": "baz",
							"lol": "kek",
						},
						MapStruct: map[string]myStruct{
							"foo": {
								Key:   "key1",
								Value: "value2",
							},
						},
						ArrayScalar: []string{"bar2", "baz"},
						ArrayStruct: []myStruct{{Value: "value1"}},
					},
				},
			},
		}

		for name := range cases {
			tc := cases[name]

			t.Run(name, func(t *testing.T) {
				tc.Init(t)

				cfg, err := config.Parse[myCfg](
					ctx,
					tc.args.flags,
					tc.args.sources...,
				)
				tc.CheckError(t, err)
				require.Equal(t, tc.results.cfg, cfg)
			})
		}
	})

	t.Run("YAML", func(t *testing.T) {
		pl_testing.Init(t)

		cases := map[string]struct {
			pl_testing.TestCase
			args    args
			results results
		}{
			"case00: no sources": {
				args: args{
					flags:   config.YAML,
					sources: []config.Source{nil, nil, nil},
				},
				results: results{},
				TestCase: pl_testing.TestCase{
					Skip:          true,
					MustFail:      true,
					MustFailAsErr: new(validator.ValidationErrors),
				},
			},
			"case01: file": {
				args: args{
					flags: config.YAML,
					sources: []config.Source{
						config.FromFile("testdata/config1.yaml"),
					},
				},
				results: results{},
				TestCase: pl_testing.TestCase{
					Skip:          true,
					MustFail:      true,
					MustFailAsErr: new(validator.ValidationErrors),
				},
			},
			"case02: file doesnt exists strict": {
				args: args{
					flags: config.YAML,
					sources: []config.Source{
						config.FromFile("testdata/foobar.yaml"),
					},
				},
				results: results{},
				TestCase: pl_testing.TestCase{
					Skip:          true,
					MustFail:      true,
					MustFailIsErr: os.ErrNotExist,
				},
			},
			"case03: file + bytes = ok": {
				args: args{
					flags: config.YAML,
					sources: []config.Source{
						config.FromFile("testdata/foobar.yaml", config.Silent),
						config.FromFile("testdata/config1.yaml"),
						config.FromFile("testdata/lolkek.yaml", config.Silent),
						config.FromBytes(config2YAML),
					},
				},
				results: results{
					cfg: myCfg{
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
						MapScalar: map[string]string{
							"foo": "bar2",
							"bar": "baz",
							"lol": "kek",
						},
						MapStruct: map[string]myStruct{
							"foo": {
								Key:   "key1",
								Value: "value2",
							},
						},
						ArrayScalar: []string{"bar2", "baz"},
						ArrayStruct: []myStruct{{Value: "value1"}},
					},
				},
			},
		}

		for name := range cases {
			tc := cases[name]

			t.Run(name, func(t *testing.T) {
				tc.Init(t)

				cfg, err := config.Parse[myCfg](
					ctx,
					tc.args.flags,
					tc.args.sources...,
				)
				tc.CheckError(t, err)
				require.Equal(t, tc.results.cfg, cfg)
			})
		}
	})
}
