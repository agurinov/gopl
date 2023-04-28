// https://github.com/edenhill/librdkafka/blob/master/CONFIGURATION.md
// https://docs.confluent.io/platform/current/clients/confluent-kafka-go/index.html
package kafka

import "fmt"

const (
	PlaintextSecurityProtocol     = "PLAINTEXT"
	SASLPlaintextSecurityProtocol = "SASL_PLAINTEXT"
	SASLSSLSecurityProtocol       = "SASL_SSL"
	SSLSecurityProtocol           = "SSL"
)

const (
	GSAPISASLMechanism       = "GSAPI"
	OAuthBearerSASLMechanism = "OAUTHBEARER"
	PlainSASLMechanism       = "PLAIN"
	SCRAMSHA256SASLMechanism = "SCRAM-SHA-256"
	SCRAMSHA512SASLMechanism = "SCRAM-SHA-512"
)

type (
	ConfigMapKey string
	ConfigMap    map[ConfigMapKey]string
)

// Common keys
const (
	BootstrapServersKey      = ConfigMapKey("bootstrap.servers")
	AllowAutoCreateTopicsKey = ConfigMapKey("allow.auto.create.topics")
	SecurityProtocolKey      = ConfigMapKey("security.protocol")
	SSLCAPEMKey              = ConfigMapKey("ssl.ca.pem")
	SASLMechanismKey         = ConfigMapKey("sasl.mechanism")
	SASLUsernameKey          = ConfigMapKey("sasl.username")
	SASLPasswordKey          = ConfigMapKey("sasl.password")
)

func (cm ConfigMap) Validate() error {
	if cm[BootstrapServersKey] == "" {
		return fmt.Errorf(
			"%w: bootstrap servers must be present",
			ErrInvalidConfigMap,
		)
	}

	return nil
}

// Consumer keys
const (
	GroupIDKey = ConfigMapKey("group.id")
)

func (cm ConfigMap) ValidateForConsumer() error {
	if err := cm.Validate(); err != nil {
		return err
	}

	if cm[GroupIDKey] == "" {
		return fmt.Errorf(
			"%w: consumer group must be present",
			ErrInvalidConfigMap,
		)
	}

	return nil
}

// Producer keys
const (
	RetryBackoffKey = ConfigMapKey("retry.backoff.ms")
	RetryCountKey   = ConfigMapKey("retries")
)

func (cm ConfigMap) ValidateForProducer() error {
	if err := cm.Validate(); err != nil {
		return err
	}

	return nil
}
