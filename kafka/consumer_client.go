package kafka

import (
	"context"
	"sync"

	"github.com/twmb/franz-go/pkg/kgo"
)

type (
	consumerClient interface {
		Close()
		Ping(context.Context) error
		CommitRecords(context.Context, ...Record) error
		PollTopicPartition(context.Context, string, int32, []int32, int) kgo.Fetches
	}
	kgoClient struct {
		client *kgo.Client
		mu     sync.RWMutex
	}
)

func (c *kgoClient) Close() {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.client == nil {
		return
	}

	c.client.Close()
}

func (c *kgoClient) Ping(ctx context.Context) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.client == nil {
		return nil
	}

	return c.client.Ping(ctx)
}

func (c *kgoClient) CommitRecords(ctx context.Context, records ...Record) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.client == nil {
		return nil
	}

	return c.client.CommitRecords(ctx, records...)
}

// NOTE: main logic appears here.
// According to kafka's documentation we CAN use pause/resume logic while we optimize goroutine
// consurrency.
// Proofs:
//  1. https://docs.confluent.io/platform/current/clients/javadocs/javadoc/org/apache/kafka/clients/consumer/KafkaConsumer.html#consumption-flow-control-heading
//  2. https://docs.confluent.io/platform/current/clients/javadocs/javadoc/org/apache/kafka/clients/consumer/KafkaConsumer.html#1-one-consumer-per-thread-heading
//
// franz-go's logic of polling separated into 2 steps:
//  1. fetch. Single eventloop across all assigned partitions with buffering feature. It is hidden and immutable for us.
//  2. poll. We can switch client to desired partition to fetch only those records. While we can't
//     fully clone the client mutex appears here only on roundtrip.
//
// There are cases when partition doesn't have records ahead and this Poll invoke can hung.
// We've limited roundtrip from config Duration and this case (deadline exceeded) for NOW interprets as IDLE state.
func (c *kgoClient) PollTopicPartition(
	ctx context.Context,
	topic string,
	targetPartition int32,
	allPartitions []int32,
	maxPollRecords int,
) kgo.Fetches {
	forPause := map[string][]int32{
		topic: allPartitions,
	}

	forResume := map[string][]int32{
		topic: {targetPartition},
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.client.PauseFetchPartitions(forPause)
	c.client.ResumeFetchPartitions(forResume)

	defer c.client.AllowRebalance()

	return c.client.PollRecords(ctx, maxPollRecords)
}

func newKgoClient(opts ...kgo.Opt) (consumerClient, error) {
	var c kgoClient

	c.mu.Lock()
	defer c.mu.Unlock()

	client, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, err
	}

	c.client = client

	return &c, nil
}
