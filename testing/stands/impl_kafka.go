package stands

import (
	"fmt"
	"net"
	"strings"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"
)

type (
	Kafka struct {
		Topics   []KafkaTopic
		Replicas int
		SASL     bool
		SSL      bool
	}
	KafkaTopic struct {
		Name       string
		Partitions int
	}
)

const (
	KafkaStandName = "kafka"
)

// https://hub.docker.com/r/bitnami/kafka
var (
	kafkaImage = docker.PullImageOptions{
		Repository: "docker.io/bitnami/kafka",
		Tag:        "3.4.0",
	}
	//nolint:gomnd
	kafkaPorts = ports{
		internal: 9093,
		external: 9094,
		cluster:  9092,
	}
)

func (s Kafka) SecurityProtocol() string {
	securityProtocol := "PLAINTEXT"

	if s.SSL {
		securityProtocol = "SSL"
	}

	if s.SASL {
		securityProtocol = "SASL_" + securityProtocol
	}

	return securityProtocol
}

func (Kafka) Name() string { return KafkaStandName }
func (s Kafka) Up(t *testing.T) bool {
	t.Helper()

	require.NotZero(t, s.Replicas)

	var (
		network      = network(t)
		kafkaCluster = newCluster(KafkaStandName, s.Replicas, kafkaPorts)

		kafka   *dockertest.Resource
		created bool
	)

	for i := range kafkaCluster {
		kafkaNode := kafkaCluster[i]

		kafka, created = container(t, &dockertest.RunOptions{
			Repository: kafkaImage.Repository,
			Tag:        kafkaImage.Tag,
			Name:       kafkaNode.Hostname(t),
			Hostname:   kafkaNode.Hostname(t),
			NetworkID:  network.ID,
			ExposedPorts: []string{
				kafkaNode.ExternalPort(),
			},
			PortBindings: map[docker.Port][]docker.PortBinding{
				docker.Port(kafkaNode.ExternalPort()): {{
					HostIP:   "localhost",
					HostPort: kafkaNode.ExternalPort(),
				}},
			},
			Env: []string{
				// TODO(a.gurinov): fmt.Sprintf("BITNAMI_DEBUG=true", tc.Debug),
				"BITNAMI_DEBUG=true",
				// https://github.com/bitnami/containers/tree/main/bitnami/kafka#configuration
				"ALLOW_PLAINTEXT_LISTENER=yes",
				"KAFKA_ENABLE_KRAFT=yes",
				fmt.Sprintf("KAFKA_KRAFT_CLUSTER_ID=%s", kafkaCluster.KafkaClusterID()),
				fmt.Sprintf("KAFKA_CFG_NODE_ID=%s", kafkaNode.KafkaNodeID()),
				fmt.Sprintf("KAFKA_CFG_CONTROLLER_QUORUM_VOTERS=%s", kafkaCluster.KafkaQuorumVoters(t)),
				"KAFKA_CFG_PROCESS_ROLES=broker,controller",
				"KAFKA_CFG_AUTO_CREATE_TOPICS_ENABLE=false",
				"KAFKA_CFG_INTER_BROKER_LISTENER_NAME=INTERNAL",
				"KAFKA_CFG_SASL_ENABLED_MECHANISMS=PLAIN",
				fmt.Sprintf("KAFKA_CFG_LISTENERS=%s",
					strings.Join([]string{
						fmt.Sprintf("CONTROLLER://:%d", kafkaPorts.cluster),
						fmt.Sprintf("INTERNAL://:%d", kafkaPorts.internal),
						fmt.Sprintf("EXTERNAL://:%s", kafkaNode.ExternalPortRaw()),
					}, ","),
				),
				fmt.Sprintf("KAFKA_CFG_ADVERTISED_LISTENERS=%s",
					strings.Join([]string{
						fmt.Sprintf("INTERNAL://%s", net.JoinHostPort(
							kafkaNode.Hostname(t),
							kafkaPorts.InternalRaw(),
						)),
						fmt.Sprintf("EXTERNAL://localhost:%s", kafkaNode.ExternalPortRaw()),
					}, ","),
				),
				fmt.Sprintf("KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=%s",
					strings.Join([]string{
						"CONTROLLER:PLAINTEXT",
						"INTERNAL:PLAINTEXT",
						"EXTERNAL:" + s.SecurityProtocol(),
					}, ","),
				),
				"KAFKA_CLIENT_USERS=sasl_user",
				"KAFKA_CLIENT_PASSWORDS=sasl_password",
			},
		})
	}

	if created {
		require.NotNil(t, kafka)

		for i := range s.Topics {
			require.NotEmpty(t, s.Topics[i].Name)
			require.NotZero(t, s.Topics[i].Partitions)

			containerExec(t, kafka, nil,
				"/opt/bitnami/kafka/bin/kafka-topics.sh",
				"--create", "--topic", s.Topics[i].Name,
				"--bootstrap-server", fmt.Sprintf("localhost:%d", kafkaPorts.internal),
				fmt.Sprintf("--replication-factor=%d", s.Replicas),
				fmt.Sprintf("--partitions=%d", s.Topics[i].Partitions),
			)
		}
	}

	return created
}
