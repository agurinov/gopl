package vault

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	vault "github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/approle"
	"github.com/hashicorp/vault/api/auth/userpass"
)

func Auth(ctx context.Context, cfg Config) (*vault.Client, error) {
	if !cfg.Enabled {
		return nil, ErrDisabled
	}

	vaultCfg := vault.DefaultConfig()
	// vaultCfg.ReadEnvironment
	vaultCfg.Address = cfg.Address

	client, err := vault.NewClient(vaultCfg)
	if err != nil {
		return nil, fmt.Errorf("vault: can't create client: %w", err)
	}

	// https://developer.hashicorp.com/vault/docs/auth/token
	if client.Token() != "" {
		return client, nil
	}

	var (
		token          string
		isAppRoleAuth  = cfg.RoleUUID != uuid.Nil
		isUserPassAuth = cfg.Username != ""
	)

	if isAppRoleAuth && isUserPassAuth {
		return nil, ErrAmbiguousAuthMethod
	}

	switch {
	case isAppRoleAuth:
		// https://developer.hashicorp.com/vault/docs/auth/approle
		authToken, authErr := authAppRole(ctx, cfg, client)
		if authErr != nil {
			return nil, fmt.Errorf("vault: can't auth via approle: %w", authErr)
		}

		token = authToken
	case isUserPassAuth:
		// https://developer.hashicorp.com/vault/docs/auth/userpass
		authToken, authErr := authUserPass(ctx, cfg, client)
		if authErr != nil {
			return nil, fmt.Errorf("vault: can't auth via username: %w", authErr)
		}

		token = authToken
	default:
		return nil, ErrUnknownAuthMethod
	}

	client.SetToken(token)

	return client, nil
}

func authAppRole(
	ctx context.Context,
	cfg Config,
	client *vault.Client,
) (string, error) {
	auth, err := approle.NewAppRoleAuth(
		cfg.RoleUUID.String(),
		&approle.SecretID{
			FromString: cfg.SecretUUID.String(),
		},
	)
	if err != nil {
		return "", err
	}

	secret, err := auth.Login(ctx, client)
	if err != nil {
		return "", err
	}

	return secret.Auth.ClientToken, nil
}

func authUserPass(
	ctx context.Context,
	cfg Config,
	client *vault.Client,
) (string, error) {
	auth, err := userpass.NewUserpassAuth(
		cfg.Username,
		&userpass.Password{
			FromString: cfg.Password,
		},
	)
	if err != nil {
		return "", err
	}

	secret, err := auth.Login(ctx, client)
	if err != nil {
		return "", err
	}

	return secret.Auth.ClientToken, nil
}
