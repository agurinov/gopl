package vault

import (
	"context"
	"fmt"

	vault "github.com/hashicorp/vault/api"

	"github.com/agurinov/gopl/env/envvars"
)

func Auth(ctx context.Context) (*vault.Client, error) {
	enabled, err := envvars.VaultEnabled.Value()

	switch {
	case err != nil:
		return nil, err
	case !enabled:
		return nil, fmt.Errorf("vault: disabled by %s", envvars.VaultEnabled)
	}

	if _, addrErr := envvars.VauldAddress.Value(); addrErr != nil {
		return nil, addrErr
	}

	cfg := vault.DefaultConfig()
	// vaultConfig.Timeout = cfg.Timeout

	client, err := vault.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("vault: can't create client: %w", err)
	}

	// https://developer.hashicorp.com/vault/docs/auth/token
	if client.Token() != "" {
		return client, nil
	}

	// https://developer.hashicorp.com/vault/docs/auth/approle
	token, err := authAppRole(ctx, client)
	if err != nil {
		return nil, err
	}

	// https://developer.hashicorp.com/vault/docs/auth/userpass
	// TODO(a.gurinov)

	client.SetToken(token)

	return client, nil
}

func authAppRole(ctx context.Context, client *vault.Client) (string, error) {
	roleID, err := envvars.VaultRoleID.Value()
	if err != nil {
		return "", fmt.Errorf("vault: can't auth via approle: %w", err)
	}

	secretID, err := envvars.VaultSecretID.Value()
	if err != nil {
		return "", fmt.Errorf("vault: can't auth via approle: %w", err)
	}

	auth, err := client.Logical().WriteWithContext(ctx, "auth/approle/login",
		map[string]any{
			"role_id":   roleID,
			"secret_id": secretID,
		},
	)
	if err != nil {
		return "", fmt.Errorf("vault: can't auth via approle: %w", err)
	}

	return auth.Auth.ClientToken, nil
}
