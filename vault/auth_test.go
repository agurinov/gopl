package vault_test

import (
	"net/url"
	"testing"

	"github.com/google/uuid"
	hashiVault "github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/require"

	pl_testing "github.com/agurinov/gopl/testing"
	"github.com/agurinov/gopl/testing/stands"
	"github.com/agurinov/gopl/vault"
)

func TestAuth(t *testing.T) {
	const (
		roleID   = "6c578acc-732c-4f66-a485-5d14a03199d3"
		secretID = "de9b68e9-99aa-4dc4-b912-816e7382698f"
		username = "dev-admin"
		password = "secret"
	)

	var (
		roleUUID   = uuid.MustParse(roleID)
		secretUUID = uuid.MustParse(secretID)
	)

	pl_testing.Init(t,
		stands.Vault{
			Replicas:   1,
			Roles:      map[uuid.UUID]uuid.UUID{roleUUID: secretUUID},
			UserPasses: map[string]string{username: password},
		},
	)

	type (
		args struct {
			cfg vault.Config
		}
		results struct{}
	)

	var (
		ctx               = t.Context()
		urlParseError     = new(url.Error)
		approleAuthError  = new(hashiVault.ResponseError)
		userpassAuthError = new(hashiVault.ResponseError)
	)

	cases := map[string]struct {
		results results
		pl_testing.TestCase
		args args
	}{
		"case00: feature disabled": {
			args: args{
				cfg: vault.Config{
					Enabled: false,
				},
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: vault.ErrDisabled,
			},
		},
		"case01: wrong address": {
			args: args{
				cfg: vault.Config{
					Enabled: true,
					Address: "%^&123456",
				},
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailAsErr: &urlParseError,
			},
		},
		"case02: unknown auth method": {
			args: args{
				cfg: vault.Config{
					Enabled: true,
					Address: "http://localhost:8200",
				},
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: vault.ErrUnknownAuthMethod,
			},
		},
		"case03: ambiguous auth method": {
			args: args{
				cfg: vault.Config{
					Enabled:  true,
					Address:  "http://localhost:8200",
					RoleUUID: roleUUID,
					Username: "username",
				},
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: vault.ErrAmbiguousAuthMethod,
			},
		},
		"case04: failed approle auth method": {
			args: args{
				cfg: vault.Config{
					Enabled:    true,
					Address:    "http://localhost:8200",
					RoleUUID:   roleUUID,
					SecretUUID: uuid.Nil,
				},
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailAsErr: &approleAuthError,
			},
		},
		"case05: success approle auth method": {
			args: args{
				cfg: vault.Config{
					Enabled:    true,
					Address:    "http://localhost:8200",
					RoleUUID:   roleUUID,
					SecretUUID: secretUUID,
				},
			},
		},
		"case04: failed userpass auth method": {
			args: args{
				cfg: vault.Config{
					Enabled:  true,
					Address:  "http://localhost:8200",
					Username: username,
					Password: "wrong password",
				},
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailAsErr: &userpassAuthError,
			},
		},
		"case05: success userpass auth method": {
			args: args{
				cfg: vault.Config{
					Enabled:  true,
					Address:  "http://localhost:8200",
					Username: username,
					Password: password,
				},
			},
		},
	}

	for name := range cases {
		tc := cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t, stands.Nop)

			client, err := vault.Auth(ctx, tc.args.cfg)
			tc.CheckError(t, err)
			require.NotNil(t, client)
			require.NotEmpty(t, client.Token())
		})
	}
}
