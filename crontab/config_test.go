package crontab_test

import (
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	"github.com/stretchr/testify/require"

	"github.com/agurinov/gopl/crontab"
	pl_testing "github.com/agurinov/gopl/testing"
)

func TestConfig_Validate(t *testing.T) {
	pl_testing.Init(t)

	type (
		args struct {
			config crontab.Config
		}
	)

	cases := map[string]struct {
		args args
		pl_testing.TestCase
	}{
		"case00: nil map": {
			args: args{
				config: crontab.Config{
					Jobs: nil,
				},
			},
			TestCase: pl_testing.TestCase{
				MustFail: true,
			},
		},
		"case01: empty map": {
			args: args{
				config: crontab.Config{
					Jobs: map[string]crontab.JobConfig{},
				},
			},
			TestCase: pl_testing.TestCase{
				MustFail: true,
			},
		},
		"case02: empty key": {
			args: args{
				config: crontab.Config{
					Jobs: map[string]crontab.JobConfig{
						"": {
							Schedule: "* * * * *",
							Timeout:  time.Second,
						},
					},
				},
			},
			TestCase: pl_testing.TestCase{
				MustFail: true,
			},
		},
		"case03: invalid job": {
			args: args{
				config: crontab.Config{
					Jobs: map[string]crontab.JobConfig{
						"foobar": {},
					},
				},
			},
			TestCase: pl_testing.TestCase{
				MustFail: true,
			},
		},
		"case04: valid job": {
			args: args{
				config: crontab.Config{
					Jobs: map[string]crontab.JobConfig{
						"foobar": {
							Schedule: "* * * * *",
							Timeout:  time.Second,
						},
					},
				},
			},
		},
	}

	for name := range cases {
		name, tc := name, cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			v := validator.New()
			require.NoError(t,
				v.RegisterValidation("notblank", validators.NotBlank),
			)

			tc.CheckError(t,
				v.Struct(tc.args.config),
			)
		})
	}
}
