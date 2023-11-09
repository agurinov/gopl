package vault

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/agurinov/gopl/env/envvars"
)

type Config struct {
	Address    string `validate:"required"`
	RoleUUID   uuid.UUID
	SecretUUID uuid.UUID
	Enabled    bool
}

// TODO(a.gurinov): must be codegen part
func LoadConfig() (Config, error) {
	var (
		enabled    bool
		address    string
		roleUUID   uuid.UUID
		secretUUID uuid.UUID
	)

	enabled, err := envvars.VaultEnabled.Value()
	if err != nil {
		return Config{}, err
	}

	address, err = envvars.VauldAddress.Value()
	if err != nil {
		return Config{}, err
	}

	roleUUID, err = envvars.VaultRoleID.Value()
	if err != nil {
		return Config{}, err
	}

	secretUUID, err = envvars.VaultSecretID.Value()
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		Enabled:    enabled,
		Address:    address,
		RoleUUID:   roleUUID,
		SecretUUID: secretUUID,
	}

	if validateErr := validator.New().Struct(cfg); validateErr != nil {
		return Config{}, validateErr
	}

	return cfg, nil
}
