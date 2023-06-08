//nolint:revive
package stands

import (
	_ "embed"
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"
)

const (
	VaultStandName = "vault"
)

var (
	// https://hub.docker.com/r/bitnami/vault
	vaultImage = docker.PullImageOptions{
		Repository: "docker.io/bitnami/vault",
		Tag:        "1.14.0",
	}
	//nolint:gomnd
	vaultPorts = ports{
		external: 8200,
		cluster:  8201,
	}
)

//go:embed etc/vault/dev-admin-policy.hcl
var vaultDevAdminPolicy string

type (
	Vault struct {
		Roles    []VaultRole
		Replicas int
	}
	VaultRole struct {
		ID       uuid.UUID
		SecretID uuid.UUID
	}
)

func (s Vault) Name() string { return VaultStandName }
func (s Vault) Up(t *testing.T) bool {
	t.Helper()

	require.Greater(t, s.Replicas, 0)

	var (
		network = network(t)
		cluster = newCluster(VaultStandName, s.Replicas, vaultPorts)

		vault   *dockertest.Resource
		created bool
	)

	for i := range cluster {
		if i > 0 {
			// TODO(a.gurinov): cluster not implemented yet
			break
		}

		node := cluster[i]

		vault, created = container(t, &dockertest.RunOptions{
			Repository: vaultImage.Repository,
			Tag:        vaultImage.Tag,
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
			CapAdd: []string{"IPC_LOCK"},
			Env: []string{
				// TODO(a.gurinov): fmt.Sprintf("BITNAMI_DEBUG=true", tc.Debug),
				"BITNAMI_DEBUG=true",
				// https://github.com/bitnami/containers/tree/main/bitnami/vault
				// https://developer.hashicorp.com/vault/tutorials/operations/configure-vault#environment-variables
				fmt.Sprintf("VAULT_ADDR=http://localhost:%d", vaultPorts.external),
				fmt.Sprintf("VAULT_API_ADDR=http://localhost:%d", vaultPorts.external),
				fmt.Sprintf("VAULT_CLUSTER_ADDR=http://localhost:%d", vaultPorts.cluster),
				fmt.Sprintf("VAULT_DEV_LISTEN_ADDRESS=0.0.0.0:%d", vaultPorts.external),
				"VAULT_TOKEN=vault-root-token",
				"VAULT_DEV_ROOT_TOKEN_ID=vault-root-token",
				"VAULT_LOG_FORMAT=json",
				"VAULT_LOG_LEVEL=debug",
			},
			Cmd: []string{
				"server",
				"-dev",
				"-dev-no-store-token",
			},
		})
	}

	if created {
		require.NotNil(t, vault)

		containerExec(t, vault, nil,
			"vault", "auth", "enable", "approle",
		)
		containerExec(t, vault, strings.NewReader(vaultDevAdminPolicy),
			"vault", "policy", "write", "dev-admin", "-",
		)

		for i := range s.Roles {
			containerExec(t, vault, nil,
				"vault", "write",
				fmt.Sprintf("auth/approle/role/role%d", i),
				fmt.Sprintf("role_id=%s", s.Roles[i].ID),
				"secret_id_num_uses=0",
				"secret_id_ttl=0",
				"token_num_uses=0",
				"token_ttl=0",
				"token_policies=dev-admin",
			)
			containerExec(t, vault, nil,
				"vault", "write",
				fmt.Sprintf("auth/approle/role/role%d/custom-secret-id", i),
				fmt.Sprintf("secret_id=%s", s.Roles[i].SecretID),
			)
		}
	}

	return created
}
