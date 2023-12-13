package stands

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"
)

type (
	Mysql struct {
		DB        string
		Liquibase string
		Replicas  int
	}
)

const (
	MysqlStandName     = "mysql"
	LiquibaseStandName = "liquibase"
)

var (
	// https://hub.docker.com/r/bitnami/mysql
	// https://github.com/bitnami/containers/blob/main/bitnami/mysql/8.2/debian-11/Dockerfile
	mysqlImage = docker.PullImageOptions{
		Repository: "docker.io/bitnami/mysql",
		Tag:        "8.2.0",
	}
	// https://hub.docker.com/r/liquibase/liquibase
	// https://github.com/liquibase/docker/blob/main/Dockerfile
	// https://github.com/liquibase/docker/blob/main/Dockerfile.alpine
	liquibaseImage = docker.PullImageOptions{
		Repository: "docker.io/liquibase/liquibase",
		Tag:        "4.25-alpine",
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
	require.NotEmpty(t, s.DB)

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
				"MYSQL_ROOT_USER=root",
				"MYSQL_ROOT_PASSWORD=root",
				fmt.Sprintf("MYSQL_DATABASE=%s", s.DB),
			},
		})
	}

	if created {
		containerExec(t, mysql, nil,
			"/opt/bitnami/mysql/bin/mysqladmin",
			"ping", "-u", "root", "--password=root", "-w",
		)
	}

	if s.Liquibase != "" {
		var (
			liquibaseNode = node{domain: LiquibaseStandName, index: 0}
			mysqlNode     = mysqlCluster[0]
		)

		liquibase, liquibaseCreated := container(t, &dockertest.RunOptions{
			Repository: liquibaseImage.Repository,
			Tag:        liquibaseImage.Tag,
			Name:       liquibaseNode.Hostname(t),
			Hostname:   liquibaseNode.Hostname(t),
			NetworkID:  network.ID,
			Entrypoint: []string{"tail", "-f", "/dev/null"},
			User:       "root:root",
			Env: []string{
				// https://docs.liquibase.com/parameters/home.html
				"INSTALL_MYSQL=true",
				"LIQUIBASE_HEADLESS=true",
				"LIQUIBASE_LOG_LEVEL=INFO",
				"LIQUIBASE_COMMAND_USERNAME=root",
				"LIQUIBASE_COMMAND_PASSWORD=root",
				"LIQUIBASE_COMMAND_CHANGELOG_FILE=migrations.sql",
				fmt.Sprintf("LIQUIBASE_COMMAND_URL=jdbc:mysql://%s:%s/%s",
					mysqlNode.Hostname(t),
					mysqlNode.ExternalPortRaw(),
					s.DB,
				),
			},
		})

		require.True(t, liquibaseCreated)

		containerExec(t, liquibase, strings.NewReader(s.Liquibase),
			"cp", "/dev/stdin", "/liquibase/changelog/migrations.sql",
		)

		containerExec(t, liquibase, nil,
			"/liquibase/docker-entrypoint.sh", "update",
		)
	}

	return created
}
