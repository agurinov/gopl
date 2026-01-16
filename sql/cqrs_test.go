package sql_test

import (
	"io"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"

	"github.com/agurinov/gopl/sql"
	pl_testing "github.com/agurinov/gopl/testing"
	"github.com/agurinov/gopl/testing/stands"
)

func TestCQRS_Connect(t *testing.T) {
	pl_testing.Init(t,
		stands.Mysql{Replicas: 1, DB: "testdb"},
	)

	type (
		args struct {
			config sql.Config
		}
	)

	ctx := t.Context()

	cases := map[string]struct {
		pl_testing.TestCase
		args args
	}{
		"case00: wrong database": {
			args: args{
				config: sql.Config{
					Driver:   "mysql",
					Host:     "localhost",
					Database: "foobar",
					User:     "root",
					Password: "root",
					Port:     3306,
				},
			},
			TestCase: pl_testing.TestCase{
				Skip:          true,
				MustFail:      true,
				MustFailIsErr: io.EOF,
			},
		},
		"case01: success": {
			args: args{
				config: sql.Config{
					Driver:   "mysql",
					Host:     "localhost",
					Database: "testdb",
					User:     "root",
					Password: "root",
					Port:     3306,
				},
			},
			TestCase: pl_testing.TestCase{
				Skip: true,
			},
		},
		"case02: noop CQRS validation error": {
			args: args{
				config: sql.Config{
					Driver:   "noop",
					Host:     "localhost",
					Database: "testdb",
					User:     "root",
					Password: "root",
					Port:     3306,
				},
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailAsErr: new(validator.ValidationErrors),
			},
		},
	}

	for name := range cases {
		tc := cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t, stands.Nop)

			err := validator.New().Struct(tc.args.config)
			require.NoError(t, err)

			db, err := sql.Connect(
				ctx,
				tc.args.config.Driver,
				tc.args.config.DSN(),
			)
			tc.CheckError(t, err)
			require.NotNil(t, db)

			cqrs, err := sql.NewCQRS(
				sql.WithRW(db),
				sql.WithRO(db),
			)
			tc.CheckError(t, err)
			require.NotNil(t, cqrs)
		})
	}
}
