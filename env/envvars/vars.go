package envvars

var GDebug = Bool("G_DEBUG")

var (
	LogEnabled = Bool("LOG_ENABLED")
	LogLevel   = String("LOG_LEVEL")
	LogFormat  = String("LOG_FORMAT")
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
	VaultEnabled          = Bool("VAULT_ENABLED")
	VauldAddress          = String("VAULT_ADDR")
	VaultToken            = String("VAULT_TOKEN")
	VaultRoleID           = UUID("VAULT_ROLE_ID")
	VaultSecretID         = UUID("VAULT_SECRET_ID")
	VaultUserpassUsername = String("VAULT_USERPASS_USERNAME")
	VaultUserpassPassword = String("VAULT_USERPASS_PASSWORD")
)

var (
	ConsulEnabled = Bool("CONSUL_ENABLED")
	ConsulAddress = String("CONSUL_HTTP_ADDR")
	ConsulToken   = String("CONSUL_HTTP_TOKEN")
)

var (
	JaegerAgentHost = String("JAEGER_AGENT_HOST")
	JaegerAgentPort = Int("JAEGER_AGENT_PORT")
)

var (
	OtelTraceEnabled  = Bool("OTEL_EXPORTER_OTLP_TRACES_ENABLED")
	OtelTraceEndpoint = String("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT")
	OtelTraceInsecure = Bool("OTEL_EXPORTER_OTLP_TRACES_INSECURE")
)

var (
	GoFile    = String("GOFILE")
	GoLine    = Int("GOLINE")
	GoPackage = String("GOPACKAGE")
)

var (
	GoMaxProcs = Int("GOMAXPROCS")
	GoMemLimit = String("GOMEMLIMIT")
)

var (
	_ = Int("GRPC_GO_LOG_VERBOSITY_LEVEL")
	_ = String("GRPC_GO_LOG_SEVERITY_LEVEL")
	_ = String("GRPC_GO_LOG_FORMATTER")
)
