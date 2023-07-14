//go:build test_unit

package segmentio_test

import (
	"context"
	"net"
	"testing"

	"github.com/agurinov/gopl/kafka"
	"github.com/agurinov/gopl/kafka/segmentio"
	pl_testing "github.com/agurinov/gopl/testing"
	"github.com/stretchr/testify/require"
)

func TestImpl_ConsumeBatch(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	tc := pl_testing.Init(t)

	var (
		ctx       = context.Background()
		configmap = kafka.ConfigMap{
			kafka.BootstrapServersKey: "localhost:9094",
			kafka.GroupIDKey:          "consumer_group_1",
		}
		kafkaCreated = tc.WithKafka(t, pl_testing.KafkaStand{
			Replicas: 1,
			SASL:     false,
			Topics: []pl_testing.KafkaTopic{
				{Name: "topic_1", Partitions: 6},
				{Name: "topic_2", Partitions: 2},
			},
		})
	)

	var (
		// 'SKIP' event means that offset must be above it
		SKIP = []byte("SKIP")
		p1o3 = []byte("p1o3")
		p2o2 = []byte("p2o2")
		p2o3 = []byte("p2o3")
		p3o1 = []byte("p3o1")
		p3o2 = []byte("p3o2")
		p3o3 = []byte("p3o3")
		p4o0 = []byte("p4o0")
		p4o1 = []byte("p4o1")
		p4o2 = []byte("p4o2")
		p4o3 = []byte("p4o3")
		p5o0 = []byte("p5o0")
		p5o1 = []byte("p5o1")
		p5o2 = []byte("p5o2")
		p5o3 = []byte("p5o3")
	)

	if kafkaCreated {
		producer := segmentio.NewProducer()
		require.NotNil(t, producer)

		require.NoError(t,
			tc.Pool(t).Retry(func() error {
				return producer.Init(ctx,
					configmap,
					kafka.Config{
						EventPosition: kafka.EventPosition{
							Topic: "topic_1",
						},
					},
				)
			}),
		)

		require.NoError(t,
			producer.ProduceBatch(ctx,
				SKIP, SKIP, SKIP, SKIP, p4o0, p5o0,
				SKIP, SKIP, SKIP, p3o1, p4o1, p5o1,
				SKIP, SKIP, p2o2, p3o2, p4o2, p5o2,
				SKIP, p1o3, p2o3, p3o3, p4o3, p5o3,
			),
		)
		require.NoError(t,
			producer.Close(),
		)
	}

	cases := map[string]struct {
		inputConfigMap        kafka.ConfigMap
		inputConfig           kafka.Config
		expectedEvents        [][]byte
		expectedEventPosition kafka.EventPosition
		pl_testing.TestCase
	}{
		"case00: partition 0": {
			inputConfigMap: configmap,
			inputConfig: kafka.Config{
				BatchSize: 2,
				EventPosition: kafka.EventPosition{
					Topic:     "topic_1",
					Partition: 0,
					Offset:    3, // lag=0
				},
			},
			expectedEvents:        [][]byte{},
			expectedEventPosition: kafka.EmptyEventPosition,
			TestCase: pl_testing.TestCase{
				Debug: true,
			},
		},
		"case01: partition 1": {
			inputConfigMap: configmap,
			inputConfig: kafka.Config{
				BatchSize: 2,
				EventPosition: kafka.EventPosition{
					Topic:     "topic_1",
					Partition: 1,
					Offset:    2, // lag=1
				},
			},
			expectedEvents:        [][]byte{p1o3},
			expectedEventPosition: kafka.EmptyEventPosition,
		},
		"case02: partition 2": {
			inputConfigMap: configmap,
			inputConfig: kafka.Config{
				BatchSize: 2,
				EventPosition: kafka.EventPosition{
					Topic:     "topic_1",
					Partition: 2,
					Offset:    1, // lag=2
				},
			},
			expectedEvents:        [][]byte{p2o2, p2o3},
			expectedEventPosition: kafka.EmptyEventPosition,
		},
		"case03: partition 3": {
			inputConfigMap: configmap,
			inputConfig: kafka.Config{
				BatchSize: 2,
				EventPosition: kafka.EventPosition{
					Topic:     "topic_1",
					Partition: 3,
					Offset:    0, // lag=3
				},
			},
			expectedEvents:        [][]byte{p3o1, p3o2},
			expectedEventPosition: kafka.EmptyEventPosition,
		},
		"case04: partition 4": {
			inputConfigMap: configmap,
			inputConfig: kafka.Config{
				BatchSize: 2,
				EventPosition: kafka.EventPosition{
					Topic:     "topic_1",
					Partition: 4,
					Offset:    -1, // lag=4
				},
			},
			expectedEvents:        [][]byte{p4o0, p4o1},
			expectedEventPosition: kafka.EmptyEventPosition,
		},
		"case05: partition 5": {
			inputConfigMap: configmap,
			inputConfig: kafka.Config{
				BatchSize: 2,
				EventPosition: kafka.EventPosition{
					Topic:     "topic_1",
					Partition: 5,
					Offset:    -1, // lag=4
				},
			},
			expectedEvents:        [][]byte{p5o0, p5o1},
			expectedEventPosition: kafka.EmptyEventPosition,
		},
	}

	for name, tc := range cases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			require.NoError(t, tc.inputConfigMap.ValidateForConsumer())
			require.NoError(t, tc.inputConfig.Validate())

			consumer := segmentio.NewConsumer()
			require.NotNil(t, consumer)
			t.Cleanup(func() { consumer.Close() })

			require.NoError(t, consumer.Init(ctx, tc.inputConfigMap, tc.inputConfig))
			require.NoError(t, consumer.Commit(ctx, tc.inputConfig.EventPosition))

			// Consume N iterations to check same returned values (cause of disabled auto commit)
			const N = 3

			for i := 0; i < N; i++ {
				events, position, err := consumer.ConsumeBatch(ctx, tc.inputConfig.BatchSize)
				tc.CheckError(t, err)
				require.Equal(t, tc.expectedEvents, events)
				require.Equal(t, tc.inputConfig.Topic, position.Topic)
			}
		})
	}
}

func testImpl_Init(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	tc := pl_testing.Init(t)

	tc.WithKafka(t, pl_testing.KafkaStand{
		SASL: false,
	})

	cases := map[string]struct {
		inputConfigMap kafka.ConfigMap
		inputConfig    kafka.Config
		pl_testing.TestCase
	}{
		"case00: wrong multiple brokers": {
			inputConfigMap: kafka.ConfigMap{
				kafka.BootstrapServersKey: "localhost:9095,localhost:9097",
				kafka.GroupIDKey:          "consumer_group_1",
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailAsErr: new(net.Error),
			},
		},
		"case01: unknown sasl mechanism": {
			inputConfigMap: kafka.ConfigMap{
				kafka.BootstrapServersKey: "localhost:9094",
				kafka.GroupIDKey:          "consumer_group_1",
				kafka.SecurityProtocolKey: "SASL_PLAINTEXT",
				kafka.SASLMechanismKey:    "foobar",
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: kafka.ErrUnknownSASLMechanism,
			},
		},
		"case02: ok SASL_PLAINTEXT": {
			inputConfigMap: kafka.ConfigMap{
				kafka.BootstrapServersKey: "localhost:9094",
				kafka.GroupIDKey:          "consumer_group_1",
				kafka.SecurityProtocolKey: "SASL_PLAINTEXT",
				kafka.SASLMechanismKey:    "PLAIN",
				kafka.SASLUsernameKey:     "sasl_user",
				kafka.SASLPasswordKey:     "sasl_password",
			},
			TestCase: pl_testing.TestCase{
				MustFail: false,
			},
		},
		"case03: invalid segmentio.ReaderConfig": {
			inputConfigMap: kafka.ConfigMap{
				kafka.BootstrapServersKey: "localhost:9094",
				kafka.GroupIDKey:          "consumer_group_1",
			},
			inputConfig: kafka.Config{
				EventPosition: kafka.EventPosition{
					Topic:     "topic_1",
					Partition: 3,
					Offset:    100500,
				},
				BatchSize: 10,
			},
			TestCase: pl_testing.TestCase{
				MustFail: true,
			},
		},
		"case03: ok PLAINTEXT": {
			inputConfigMap: kafka.ConfigMap{
				kafka.BootstrapServersKey: "localhost:9094",
				kafka.GroupIDKey:          "consumer_group_1",
			},
			inputConfig: kafka.Config{
				EventPosition: kafka.EventPosition{
					Topic: "topic_1",
				},
				BatchSize: 10,
			},
			TestCase: pl_testing.TestCase{
				MustFail: false,
				Debug:    true,
			},
		},
	}

	for name, tc := range cases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			var (
				ctx  = context.Background()
				impl = segmentio.New()
			)

			tc.CheckError(t,
				tc.Pool(t).Retry(func() error {
					return impl.Init(
						ctx,
						tc.inputConfigMap,
						tc.inputConfig,
					)
				}),
			)
			t.Cleanup(func() { impl.Close() })
		})
	}
}
