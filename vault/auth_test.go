package vault_test

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/agurinov/gopl/env/envvars"
	pl_testing "github.com/agurinov/gopl/testing"
	"github.com/agurinov/gopl/testing/stands"
	"github.com/agurinov/gopl/vault"
)

func TestAuth(t *testing.T) {
	const (
		roleID   = "6c578acc-732c-4f66-a485-5d14a03199d3"
		secretID = "de9b68e9-99aa-4dc4-b912-816e7382698f"
	)

	pl_testing.Init(t,
		stands.Vault{
			Replicas: 1,
			Roles: []stands.VaultRole{
				{
					ID:       uuid.MustParse(roleID),
					SecretID: uuid.MustParse(secretID),
				},
			},
		},
	)

	os.Setenv(envvars.VaultEnabled.String(), "true")
	os.Setenv(envvars.VauldAddress.String(), "http://localhost:8200")
	os.Setenv(envvars.VaultRoleID.String(), roleID)
	os.Setenv(envvars.VaultSecretID.String(), secretID)

	ctx := context.TODO()

	cfg, err := vault.LoadConfig()
	require.NoError(t, err)

	client, err := vault.Auth(ctx, cfg)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.NotEmpty(t, client.Token())
}
