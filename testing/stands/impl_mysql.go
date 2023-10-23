package stands

import (
	"fmt"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"
)

type (
	Mysql struct {
		DB       string
		Replicas int
	}
)

const (
	MysqlStandName = "mysql"
)

var (
	// https://hub.docker.com/r/bitnami/mysql
	mysqlImage = docker.PullImageOptions{
		Repository: "docker.io/bitnami/mysql",
		Tag:        "8.0.33",
	}
	//nolint:gomnd
	mysqlPorts = ports{
		external: 3306,
	}
)

func (Mysql) Name() string { return MysqlStandName }
func (s Mysql) Up(t *testing.T) bool {
	t.Helper()

	require.NotZero(t, s.Replicas)

	var (
		network      = network(t)
		mysqlCluster = newCluster(MysqlStandName, s.Replicas, mysqlPorts)

		mysql   *dockertest.Resource
		created bool
	)

	for i := range mysqlCluster {
		if i > 0 {
			// TODO(a.gurinov): cluster not implemented yet
			break
		}

		mysqlNode := mysqlCluster[i]

		mysql, created = container(t, &dockertest.RunOptions{
			Repository: mysqlImage.Repository,
			Tag:        mysqlImage.Tag,
			Name:       mysqlNode.Hostname(t),
			Hostname:   mysqlNode.Hostname(t),
			NetworkID:  network.ID,
			ExposedPorts: []string{
				mysqlNode.ExternalPort(),
			},
			PortBindings: map[docker.Port][]docker.PortBinding{
				docker.Port(mysqlNode.ExternalPort()): {{
					HostIP:   "localhost",
					HostPort: mysqlNode.ExternalPort(),
				}},
			},
			Env: []string{
				// TODO(a.gurinov): fmt.Sprintf("BITNAMI_DEBUG=true", tc.Debug),
				"BITNAMI_DEBUG=true",
				// https://github.com/bitnami/containers/tree/main/bitnami/mysql
				"ALLOW_EMPTY_PASSWORD=yes",
				fmt.Sprintf("MYSQL_DATABASE=%s", s.DB),
			},
		})
	}

	if created {
		require.NotNil(t, mysql)
	}

	return created
}
