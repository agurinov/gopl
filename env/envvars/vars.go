package pl_envvars

var (
	LogLevel   = String("LOG_LEVEL")
	LogDriver  = String("LOG_DRIVER")
	LogEnabled = Bool("LOG_ENABLED")
)

var (
	DB_HOST           = IP("DB_HOST")
	DB_REPLICATION_ID = UUID("DB_REPLICATION_ID")
)
