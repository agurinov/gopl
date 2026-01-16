package sql_test

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"

	"github.com/agurinov/gopl/sql"
	pl_testing "github.com/agurinov/gopl/testing"
)

var (
	//go:embed testdata/oneline.sql
	oneline string

	//go:embed testdata/multiline.sql
	multiline string
)

func TestQuery_WithSpan(t *testing.T) {
	pl_testing.Init(t)

	type (
		args struct {
			query sql.Query
		}
		results struct {
			stmt string
		}
	)

	ctx := t.Context()

	cases := map[string]struct {
		args    args
		results results
		pl_testing.TestCase
	}{
		"case00: oneline": {
			args: args{
				query: sql.Query(oneline),
			},
			results: results{
				stmt: `select true;
`,
			},
		},
		"case01: multiline": {
			args: args{
				query: sql.Query(multiline),
			},
			results: results{
				stmt: `select
	true as one
	false as two
;
`,
			},
		},
	}

	for name := range cases {
		tc := cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			span := trace.SpanFromContext(ctx)

			stmt := tc.args.query.WithSpan(span)
			require.Equal(t, tc.results.stmt, stmt)
		})
	}
}
