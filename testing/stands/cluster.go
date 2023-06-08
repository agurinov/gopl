package stands

import (
	"encoding/base64"
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
)

type (
	ports struct {
		internal int
		external int
		cluster  int
	}
	node struct {
		domain string
		ports  ports
		index  int
	}
	cluster []node
)

func (p ports) Cluster() string               { return fmt.Sprintf("%d/tcp", p.cluster) }
func (p ports) Internal() string              { return fmt.Sprintf("%d/tcp", p.internal) }
func (p ports) InternalRaw() string           { return fmt.Sprintf("%d", p.internal) }
func (p ports) External(offset int) string    { return fmt.Sprintf("%d/tcp", p.external+offset) }
func (p ports) ExternalRaw(offset int) string { return fmt.Sprintf("%d", p.external+offset) }

func (n node) Hostname(t *testing.T) string {
	t.Helper()

	return fmt.Sprintf("%s_%s%d", hash(t), n.domain, n.index)
}
func (n node) ExternalPort() string    { return n.ports.External(n.index) }
func (n node) ExternalPortRaw() string { return n.ports.ExternalRaw(n.index) }
func (n node) KafkaNodeID() string     { return fmt.Sprintf("%d", n.index) }

func (c cluster) KafkaQuorumVoters(t *testing.T) string {
	t.Helper()

	s := make([]string, 0, len(c))

	for i := range c {
		s = append(s,
			fmt.Sprintf("%d@%s:%d", c[i].index, c[i].Hostname(t), c[i].ports.cluster),
		)
	}

	return strings.Join(s, ",")
}

func (c cluster) KafkaClusterID() string {
	clusterID := uuid.Nil

	if b, err := clusterID.MarshalBinary(); err == nil {
		return base64.StdEncoding.EncodeToString(b)
	}

	return ""
}

func newCluster(domain string, replicas int, ports ports) cluster {
	c := make(cluster, 0, replicas)

	for i := 0; i < replicas; i++ {
		c = append(c, node{
			index:  i,
			domain: domain,
			ports:  ports,
		})
	}

	return c
}
