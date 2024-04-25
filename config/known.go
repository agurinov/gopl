package config

import (
	"github.com/agurinov/gopl/diag/log"
	"github.com/agurinov/gopl/diag/probes"
	"github.com/agurinov/gopl/graceful"
	"github.com/agurinov/gopl/grpc"
	"github.com/agurinov/gopl/http"
	"github.com/agurinov/gopl/vault"
)

type (
	Vault    = vault.Config
	Logger   = log.Config
	Probes   = probes.Config
	Graceful = graceful.Config
	GRPC     = grpc.Config
	HTTP     = http.Config
)
