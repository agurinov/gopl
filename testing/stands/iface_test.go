package stands_test

import (
	_ "embed"
	"testing"

	pl_testing "github.com/agurinov/gopl/testing"
	"github.com/agurinov/gopl/testing/stands"
)

//go:embed etc/mysql/liquibase/migrations.sql
var liquibase string

func TestStandsUp(t *testing.T) {
	pl_testing.Init(t,
		stands.Kafka{
			Replicas: 3,
			Topics: []stands.KafkaTopic{
				{Name: "topic_1", Partitions: 3},
				{Name: "topic_2", Partitions: 6},
			},
		},
		stands.Mysql{
			Replicas:  1,
			DB:        "testdb",
			Liquibase: liquibase,
		},
	)
}
