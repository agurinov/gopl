package stands

import (
	"fmt"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"
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

type (
	Mysql struct {
		DB       string
		Replicas int
	}
)

func (s Mysql) Name() string { return MysqlStandName }
func (s Mysql) Up(t *testing.T) bool {
	t.Helper()

	t.Helper()

	require.Greater(t, s.Replicas, 0)

	var (
		network = network(t)
		cluster = newCluster(MysqlStandName, s.Replicas, mysqlPorts)

		mysql   *dockertest.Resource
		created bool
	)

	for i := range cluster {
		if i > 0 {
			// TODO(a.gurinov): cluster not implemented yet
			break
		}

		node := cluster[i]

		mysql, created = container(t, &dockertest.RunOptions{
			Repository: mysqlImage.Repository,
			Tag:        mysqlImage.Tag,
			Name:       node.Hostname(t),
			Hostname:   node.Hostname(t),
			NetworkID:  network.ID,
			ExposedPorts: []string{
				node.ExternalPort(),
			},
			PortBindings: map[docker.Port][]docker.PortBinding{
				docker.Port(node.ExternalPort()): {{
					HostIP:   "localhost",
					HostPort: node.ExternalPort(),
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
