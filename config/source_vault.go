package config

import (
	"context"

	vault "github.com/hashicorp/vault/api"
)

func FromVault(
	_ *vault.Client,
	_ string,
) Source {
	return func(context.Context) ([]byte, error) {
		return nil, nil
	}
}
