package vault_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"

	pl_testing "github.com/agurinov/gopl/testing"
	"github.com/agurinov/gopl/vault"
)

func TestLoadConfig(t *testing.T) {
	pl_testing.Init(t)

	type (
		args    struct{}
		results struct {
			cfg vault.Config
		}
	)

	require.NoError(t, godotenv.Load("testdata/.env"))

	cases := map[string]struct {
		args    args
		results results
		pl_testing.TestCase
	}{
		"case00: success godotenv": {
			results: results{
				cfg: vault.Config{
					Enabled:    true,
					Address:    "http://localhost:8200",
					RoleUUID:   uuid.MustParse("6c578acc-732c-4f66-a485-5d14a03199d3"),
					SecretUUID: uuid.MustParse("de9b68e9-99aa-4dc4-b912-816e7382698f"),
					Username:   "dev-tester",
					Password:   "dev-tester-password",
				},
			},
		},
	}

	for name := range cases {
		tc := cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			cfg, err := vault.LoadConfig()
			tc.CheckError(t, err)
			require.Equal(t, tc.results.cfg, cfg)
		})
	}
}
