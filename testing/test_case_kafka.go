package pl_testing

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"
)

type (
	KafkaStand struct {
		Topics    []KafkaTopic
		SASL      bool
		SSL       bool
		ZooKeeper bool
	}
	KafkaTopic struct {
		Name       string
		Partitions uint64
	}
)

func (s KafkaStand) SecurityProtocol() string {
	securityProtocol := "PLAINTEXT"

	if s.SSL {
		securityProtocol = "SSL"
	}

	if s.SASL {
		securityProtocol = "SASL_" + securityProtocol
	}

	return securityProtocol
}

func (tc TestCase) WithKafka(t *testing.T, opts KafkaStand) bool {
	t.Helper()

	var (
		network           = tc.network(t, "gopl_kafka_network")
		kafkaZooKeeperEnv []string
	)

	if opts.ZooKeeper {
		tc.container(t, &dockertest.RunOptions{
			Name:       "zookeeper",
			Repository: "bitnami/zookeeper",
			Tag:        "latest",
			Hostname:   "zookeeper",
			NetworkID:  network.ID,
			Env: []string{
				fmt.Sprintf("BITNAMI_DEBUG=%t", tc.Debug),
				"ALLOW_ANONYMOUS_LOGIN=yes",
			},
		})

		kafkaZooKeeperEnv = []string{
			"KAFKA_CFG_ZOOKEEPER_CONNECT=zookeeper:2181",
			"KAFKA_CFG_ZOOKEEPER_PROTOCOL=PLAINTEXT",
		}
	}

	kafka := tc.container(t, &dockertest.RunOptions{
		Name:       "kafka",
		Repository: "bitnami/kafka",
		Tag:        "3.1.0",
		Hostname:   "kafka",
		NetworkID:  network.ID,
		ExposedPorts: []string{
			"9094/tcp",
		},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"9094/tcp": {{HostIP: "localhost", HostPort: "9094/tcp"}},
		},
		Env: append([]string{
			fmt.Sprintf("BITNAMI_DEBUG=%t", tc.Debug),
			"KAFKA_BROKER_ID=1",
			"ALLOW_PLAINTEXT_LISTENER=yes",
			"KAFKA_CFG_AUTO_CREATE_TOPICS_ENABLE=false",
			"KAFKA_CFG_INTER_BROKER_LISTENER_NAME=INTERNAL",
			"KAFKA_CFG_LISTENERS=INTERNAL://:9092,EXTERNAL://:9094",
			"KAFKA_CFG_ADVERTISED_LISTENERS=INTERNAL://kafka:9092,EXTERNAL://localhost:9094",
			fmt.Sprintf(
				"KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=INTERNAL:PLAINTEXT,EXTERNAL:%s",
				opts.SecurityProtocol(),
			),
			"KAFKA_CLIENT_USERS=sasl_user",
			"KAFKA_CLIENT_PASSWORDS=sasl_password",
		}, kafkaZooKeeperEnv...),
	})

	created := kafka != nil

	if created {
		for i := range opts.Topics {
			exitCode, err := kafka.Exec([]string{
				"/opt/bitnami/kafka/bin/kafka-topics.sh",
				"--bootstrap-server", "localhost:9092",
				"--create", "--topic", opts.Topics[i].Name,
				"--replication-factor", "1",
				"--partitions", strconv.FormatUint(opts.Topics[i].Partitions, 10),
			}, dockertest.ExecOptions{})

			require.NoError(t, err)
			require.Equal(t, 0, exitCode)
		}
	}

	return created
}
