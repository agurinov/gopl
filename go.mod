module github.com/agurinov/gopl

go 1.21

toolchain go1.21.10

require (
	github.com/bufbuild/protovalidate-go v0.2.1
	github.com/dustin/go-humanize v1.0.1
	github.com/go-chi/chi/v5 v5.0.11
	github.com/go-playground/validator/v10 v10.14.1
	github.com/go-sql-driver/mysql v1.8.1
	github.com/goccy/go-yaml v1.11.3
	github.com/google/uuid v1.6.0
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.0.1
	github.com/hashicorp/vault/api v1.10.0
	github.com/hashicorp/vault/api/auth/approle v0.5.0
	github.com/hashicorp/vault/api/auth/userpass v0.5.0
	github.com/improbable-eng/grpc-web v0.15.0
	github.com/jackc/pgx/v5 v5.6.0
	github.com/jmoiron/sqlx v1.4.0
	github.com/joho/godotenv v1.5.1
	github.com/ory/dockertest/v3 v3.10.0
	github.com/prometheus/client_golang v1.18.0
	github.com/stretchr/testify v1.9.0
	go.opentelemetry.io/otel v1.27.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.27.0
	go.opentelemetry.io/otel/sdk v1.27.0
	go.opentelemetry.io/otel/trace v1.27.0
	go.uber.org/automaxprocs v1.5.3
	go.uber.org/mock v0.4.0
	go.uber.org/zap v1.26.0
	golang.org/x/exp v0.0.0-20230522175609-2e198f4a06a1
	golang.org/x/sync v0.7.0
	golang.org/x/tools v0.22.0
	google.golang.org/grpc v1.64.0
	google.golang.org/protobuf v1.34.1
	gopkg.in/yaml.v3 v3.0.1
)

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.31.0-20230802163732-1c33ebd9ecfa.1 // indirect
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78 // indirect
	github.com/Microsoft/go-winio v0.6.0 // indirect
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5 // indirect
	github.com/antlr/antlr4/runtime/Go/antlr/v4 v4.0.0-20230512164433-5d1fd1a340c9 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cenkalti/backoff/v3 v3.0.0 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/containerd/continuity v0.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/desertbit/timer v0.0.0-20180107155436-c41aec40b27f // indirect
	github.com/docker/cli v20.10.17+incompatible // indirect
	github.com/docker/docker v20.10.7+incompatible // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/fatih/color v1.10.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.2 // indirect
	github.com/go-jose/go-jose/v3 v3.0.1 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/cel-go v0.17.1 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.20.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-retryablehttp v0.6.6 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/go-secure-stdlib/parseutil v0.1.6 // indirect
	github.com/hashicorp/go-secure-stdlib/strutil v0.1.2 // indirect
	github.com/hashicorp/go-sockaddr v1.0.2 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/klauspost/compress v1.11.7 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/matttproud/golang_protobuf_extensions/v2 v2.0.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/moby/term v0.0.0-20201216013528-df9cb8a40635 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.2 // indirect
	github.com/opencontainers/runc v1.1.5 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.45.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	github.com/rs/cors v1.7.0 // indirect
	github.com/ryanuber/go-glob v1.0.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/stoewer/go-strcase v1.3.0 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20180127040702-4e3ac2762d5f // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.2.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.27.0 // indirect
	go.opentelemetry.io/otel/metric v1.27.0 // indirect
	go.opentelemetry.io/proto/otlp v1.2.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/crypto v0.24.0 // indirect
	golang.org/x/mod v0.18.0 // indirect
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	golang.org/x/time v0.0.0-20200416051211-89c76fbcd5d1 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240520151616-dc85e6b867a5 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240515191416-fc5f0ca64291 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	nhooyr.io/websocket v1.8.6 // indirect
)
