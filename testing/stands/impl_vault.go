//revive:disable:add-constant
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

type (
	Vault struct {
		Roles      map[uuid.UUID]uuid.UUID
		UserPasses map[string]string
		Replicas   int
	}
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
	//nolint:gomnd,mnd
	vaultPorts = ports{
		external: 8200,
		cluster:  8201,
	}
)

//go:embed etc/vault/dev-admin-policy.hcl
var vaultDevAdminPolicy string

func (Vault) Name() string { return VaultStandName }

//nolint:perfsprint
func (s Vault) Up(t *testing.T) bool {
	t.Helper()

	require.NotZero(t, s.Replicas)

	var (
		network      = network(t)
		vaultCluster = newCluster(VaultStandName, s.Replicas, vaultPorts)

		vault   *dockertest.Resource
		created bool
	)

	for i := range vaultCluster {
		if i > 0 {
			// TODO(a.gurinov): cluster not implemented yet
			break
		}

		vaultNode := vaultCluster[i]

		vault, created = container(t, &dockertest.RunOptions{
			Repository: vaultImage.Repository,
			Tag:        vaultImage.Tag,
			Name:       vaultNode.Hostname(t),
			Hostname:   vaultNode.Hostname(t),
			NetworkID:  network.ID,
			ExposedPorts: []string{
				vaultNode.ExternalPort(),
			},
			PortBindings: map[docker.Port][]docker.PortBinding{
				docker.Port(vaultNode.ExternalPort()): {{
					HostIP:   "localhost",
					HostPort: vaultNode.ExternalPort(),
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
		containerExec(t, vault, strings.NewReader(vaultDevAdminPolicy),
			"vault", "policy", "write", "dev-admin", "-",
		)
		containerExec(t, vault, nil,
			"vault", "auth", "enable", "approle",
		)
		containerExec(t, vault, nil,
			"vault", "auth", "enable", "userpass",
		)

		for roleUUID, secretUUID := range s.Roles {
			require.NotEqual(t, roleUUID, uuid.Nil)
			require.NotEqual(t, secretUUID, uuid.Nil)

			containerExec(t, vault, nil,
				"vault", "write",
				fmt.Sprintf("auth/approle/role/%s", roleUUID),
				fmt.Sprintf("role_id=%s", roleUUID),
				"secret_id_num_uses=0",
				"secret_id_ttl=0",
				"token_num_uses=0",
				"token_ttl=0",
				"token_policies=dev-admin",
			)
			containerExec(t, vault, nil,
				"vault", "write",
				fmt.Sprintf("auth/approle/role/%s/custom-secret-id", roleUUID),
				fmt.Sprintf("secret_id=%s", secretUUID),
			)
		}

		for username, password := range s.UserPasses {
			require.NotEmpty(t, username)
			require.NotEmpty(t, password)

			containerExec(t, vault, nil,
				"vault", "write",
				fmt.Sprintf("auth/userpass/users/%s", username),
				fmt.Sprintf("password=%s", password),
				"policies=dev-admin",
			)
		}
	}

	return created
}
