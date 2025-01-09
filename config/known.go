package config

import (
	"github.com/agurinov/gopl/crontab"
	"github.com/agurinov/gopl/diag/log"
	"github.com/agurinov/gopl/diag/metrics"
	"github.com/agurinov/gopl/diag/probes"
	"github.com/agurinov/gopl/diag/trace"
	"github.com/agurinov/gopl/graceful"
	"github.com/agurinov/gopl/grpc"
	"github.com/agurinov/gopl/http"
	"github.com/agurinov/gopl/sql"
	"github.com/agurinov/gopl/telegram"
	"github.com/agurinov/gopl/vault"
)

type (
	Logger  = log.Config
	Probes  = probes.Config
	Trace   = trace.Config
	Metrics = metrics.Config
)

type (
	Vault    = vault.Config
	SQL      = sql.Config
	Graceful = graceful.Config
	GRPC     = grpc.Config
	HTTP     = http.Config
)

type (
	Crontab = crontab.Config
)

type (
	TelegramWebApp = telegram.WebAppConfig
	TelegramBot    = telegram.BotConfig
)
