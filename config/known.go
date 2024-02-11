package config

import (
	"github.com/agurinov/gopl/diag/log"
	"github.com/agurinov/gopl/graceful"
	"github.com/agurinov/gopl/vault"
)

type (
	Vault    = vault.Config
	Logger   = log.Config
	Graceful = graceful.Config
)
