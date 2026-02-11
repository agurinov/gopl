package kafka

import (
	"github.com/twmb/franz-go/pkg/kgo"
)

func (c consumer[R, V]) kgoOptions() []kgo.Opt {
	var (
		configOpts   = c.config.kgoOptions()
		consumerOpts = []kgo.Opt{
			kgo.DisableAutoCommit(),
			kgo.BlockRebalanceOnPoll(),
			kgo.OnPartitionsAssigned(c.onAssigned),
			kgo.OnPartitionsRevoked(c.onRevoked),
			kgo.OnPartitionsLost(c.onRevoked),
		}
		// TODO: custom options to override all
	)

	return append(configOpts, consumerOpts...)
}
