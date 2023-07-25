package envvars

var GDebug = Bool("G_DEBUG")

var (
	LogEnabled = Bool("LOG_ENABLED")
	LogLevel   = String("LOG_LEVEL")
	LogDriver  = String("LOG_DRIVER")
)

var (
	DBEnabled       = Bool("DB_ENABLED")
	DBHost          = IP("DB_HOST")
	DBReplicationID = UUID("DB_REPLICATION_ID")
)

var (
	KafkaEnabled       = Bool("KFK_ENABLED")
	KafkaBrokerURL     = URL("KFK_BROKER_URL")
	KafkaTopic         = String("KFK_TOPIC")
	KafkaConsumerGroup = String("KFK_CONSUMER_GROUP")
	KafkaOffset        = Int("KFK_OFFSET")
	KafkaReadTimeout   = Duration("KFK_READ_TIMEOUT")
)

var (
	VaultEnabled  = Bool("VAULT_ENABLED")
	VauldAddress  = String("VAULT_ADDR")
	VaultToken    = String("VAULT_TOKEN")
	VaultRoleID   = UUID("VAULT_ROLE_ID")
	VaultSecretID = UUID("VAULT_SECRET_ID")
)
