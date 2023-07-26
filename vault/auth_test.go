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
	pl_testing.Init(t,
		stands.Vault{
			Replicas: 1,
			Roles: []stands.VaultRole{
				{
					ID:       uuid.MustParse("6c578acc-732c-4f66-a485-5d14a03199d3"),
					SecretID: uuid.MustParse("de9b68e9-99aa-4dc4-b912-816e7382698f"),
				},
			},
		},
	)

	os.Setenv(envvars.VaultEnabled.String(), "true")
	os.Setenv(envvars.VauldAddress.String(), "http://localhost:8200")
	os.Setenv(envvars.VaultRoleID.String(), "6c578acc-732c-4f66-a485-5d14a03199d3")
	os.Setenv(envvars.VaultSecretID.String(), "de9b68e9-99aa-4dc4-b912-816e7382698f")

	ctx := context.TODO()

	client, err := vault.Auth(ctx)
	require.NoError(t, err)
	require.NotNil(t, client)
}
