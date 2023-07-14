package pl_testing

import (
	"encoding/base64"
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"
)

const (
	KafkaStandName StandName = "kafka"
	imageRepo                = "docker.io/bitnami/kafka"
	imageTag                 = "3.4.0"
	controllerPort           = 9092
	internalPort             = 9093
	externalPort             = 9094
)

type (
	KafkaStand struct {
		Replicas int
		Topics   []KafkaTopic
		SASL     bool
		SSL      bool
	}
	KafkaTopic struct {
		Name       string
		Partitions int
	}
	node struct {
		hostname string
		quorum   string
		port     int
	}
	cluster []node
)

func (n node) ExternalPort() string {
	return fmt.Sprintf("%d/tcp", n.port)
}

func (c cluster) Quorum() string {
	s := make([]string, 0, len(c))

	for i := range c {
		s = append(s, c[i].quorum)
	}

	return strings.Join(s, ",")
}

func (c cluster) ID() string {
	clusterID := uuid.Nil

	if b, err := clusterID.MarshalBinary(); err == nil {
		return base64.StdEncoding.EncodeToString(b)
	}

	return ""
}

func (c *cluster) fill(replicas int) {
	for i := 0; i < replicas; i++ {
		hostname := fmt.Sprintf("kafka%d", i)

		*c = append(*c, node{
			hostname: hostname,
			quorum:   fmt.Sprintf("%d@%s:%d", i, hostname, controllerPort),
			port:     externalPort + i,
		})
	}
}

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

func (s KafkaStand) Name() StandName { return KafkaStandName }
func (s KafkaStand) Up(t *testing.T) bool {
	t.Helper()

	return false
}

func (tc TestCase) WithKafka(t *testing.T, opts KafkaStand) bool {
	t.Helper()

	require.Greater(t, opts.Replicas, 0)

	var (
		network = tc.network(t, "gopl_kafka_network")
		cluster = make(cluster, 0, opts.Replicas)
	)

	cluster.fill(opts.Replicas)

	var (
		kafka   *dockertest.Resource
		created bool
	)

	for i, node := range cluster {
		kafka, created = tc.container(t, &dockertest.RunOptions{
			Repository: imageRepo,
			Tag:        imageTag,
			Name:       node.hostname,
			Hostname:   node.hostname,
			NetworkID:  network.ID,
			ExposedPorts: []string{
				node.ExternalPort(),
			},
			PortBindings: map[docker.Port][]docker.PortBinding{
				docker.Port(node.ExternalPort()): {{
					HostIP:   "localhost",
					HostPort: node.ExternalPort(),
				}},
			},
			Env: []string{
				fmt.Sprintf("BITNAMI_DEBUG=%t", tc.Debug),
				// https://github.com/bitnami/containers/tree/main/bitnami/kafka#configuration
				"ALLOW_PLAINTEXT_LISTENER=yes",
				"KAFKA_ENABLE_KRAFT=yes",
				fmt.Sprintf("KAFKA_KRAFT_CLUSTER_ID=%s", cluster.ID()),
				fmt.Sprintf("KAFKA_CFG_NODE_ID=%d", i),
				fmt.Sprintf("KAFKA_CFG_CONTROLLER_QUORUM_VOTERS=%s", cluster.Quorum()),
				"KAFKA_CFG_PROCESS_ROLES=broker,controller",
				"KAFKA_CFG_AUTO_CREATE_TOPICS_ENABLE=false",
				"KAFKA_CFG_INTER_BROKER_LISTENER_NAME=INTERNAL",
				"KAFKA_CFG_SASL_ENABLED_MECHANISMS=PLAIN",
				fmt.Sprintf("KAFKA_CFG_LISTENERS=%s",
					strings.Join([]string{
						fmt.Sprintf("CONTROLLER://:%d", controllerPort),
						fmt.Sprintf("INTERNAL://:%d", internalPort),
						fmt.Sprintf("EXTERNAL://:%d", node.port),
					}, ","),
				),
				fmt.Sprintf("KAFKA_CFG_ADVERTISED_LISTENERS=%s",
					strings.Join([]string{
						fmt.Sprintf("INTERNAL://%s:%d", node.hostname, internalPort),
						fmt.Sprintf("EXTERNAL://localhost:%d", node.port),
					}, ","),
				),
				fmt.Sprintf("KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=%s",
					strings.Join([]string{
						"CONTROLLER:PLAINTEXT",
						"INTERNAL:PLAINTEXT",
						"EXTERNAL:" + opts.SecurityProtocol(),
					}, ","),
				),
				"KAFKA_CLIENT_USERS=sasl_user",
				"KAFKA_CLIENT_PASSWORDS=sasl_password",
			},
		})
	}

	if created {
		for i := range opts.Topics {
			exitCode, err := kafka.Exec([]string{
				"/opt/bitnami/kafka/bin/kafka-topics.sh",
				"--create", "--topic", opts.Topics[i].Name,
				"--bootstrap-server", fmt.Sprintf("localhost:%d", internalPort),
				fmt.Sprintf("--replication-factor=%d", opts.Replicas),
				fmt.Sprintf("--partitions=%d", opts.Topics[i].Partitions),
			}, dockertest.ExecOptions{})

			require.NoError(t, err)
			require.Equal(t, 0, exitCode)
		}
	}

	return created
}
