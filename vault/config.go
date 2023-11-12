package vault

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/agurinov/gopl/env/envvars"
)

type Config struct {
	Address    string `validate:"required"`
	Username   string
	Password   string
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
		username   string
		password   string
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

	username, err = envvars.VaultUserpassUsername.Value()
	if err != nil {
		return Config{}, err
	}

	password, err = envvars.VaultUserpassPassword.Value()
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		Enabled:    enabled,
		Address:    address,
		RoleUUID:   roleUUID,
		SecretUUID: secretUUID,
		Username:   username,
		Password:   password,
	}

	if validateErr := validator.New().Struct(cfg); validateErr != nil {
		return Config{}, validateErr
	}

	return cfg, nil
}
