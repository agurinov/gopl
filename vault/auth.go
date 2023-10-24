package vault

import (
	"context"
	"fmt"

	vault "github.com/hashicorp/vault/api"
)

func Auth(ctx context.Context, cfg Config) (*vault.Client, error) {
	if !cfg.Enabled {
		return nil, fmt.Errorf("vault: feature disabled")
	}

	vaultCfg := vault.DefaultConfig()
	vaultCfg.Address = cfg.Address

	client, err := vault.NewClient(vaultCfg)
	if err != nil {
		return nil, fmt.Errorf("vault: can't create client: %w", err)
	}

	// https://developer.hashicorp.com/vault/docs/auth/token
	if client.Token() != "" {
		return client, nil
	}

	// https://developer.hashicorp.com/vault/docs/auth/approle
	token, err := authAppRole(ctx, client, cfg)
	if err != nil {
		return nil, err
	}

	// https://developer.hashicorp.com/vault/docs/auth/userpass
	// TODO(a.gurinov)

	client.SetToken(token)

	return client, nil
}

func authAppRole(
	ctx context.Context,
	client *vault.Client,
	cfg Config,
) (string, error) {
	auth, err := client.Logical().WriteWithContext(ctx,
		"auth/approle/login",
		map[string]any{
			"role_id":   cfg.RoleUUID,
			"secret_id": cfg.SecretUUID,
		},
	)
	if err != nil {
		return "", fmt.Errorf("vault: can't auth via approle: %w", err)
	}

	return auth.Auth.ClientToken, nil
}
