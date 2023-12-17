package vault_test

import (
	_ "embed"
	"testing"

	pl_testing "github.com/agurinov/gopl/testing"
	"github.com/agurinov/gopl/testing/stands"
)

//go:embed testdata/migrations.sql
var liquibase string

func TestLiquibase_MigrationUp(t *testing.T) {
	pl_testing.Init(t,
		stands.Mysql{
			Replicas:  1,
			DB:        "testdb",
			Liquibase: liquibase,
		},
	)
}
