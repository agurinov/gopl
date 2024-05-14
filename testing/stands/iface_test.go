package stands_test

import (
	"embed"
	"io/fs"
	"testing"

	"github.com/stretchr/testify/require"

	pl_testing "github.com/agurinov/gopl/testing"
	"github.com/agurinov/gopl/testing/stands"
)

//go:embed etc/mysql/liquibase
var liquibase embed.FS

func TestStandsUp(t *testing.T) {
	trimmed, err := fs.Sub(liquibase, "etc/mysql/liquibase/tree")
	require.NoError(t, err)

	pl_testing.Init(t,
		stands.Kafka{
			Replicas: 3,
			Topics: []stands.KafkaTopic{
				{Name: "topic_1", Partitions: 3},
				{Name: "topic_2", Partitions: 6},
			},
		},
		stands.Mysql{
			Replicas: 1,
			DB:       "testdb",
			Liquibase: stands.Liquibase{
				Enabled:    true,
				FS:         trimmed,
				Entrypoint: "changelog.xml",
			},
		},
	)
}
