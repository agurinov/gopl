//nolint:testableexamples
package kafka_test

import (
	"context"
	"time"

	"github.com/agurinov/gopl/backoff"
	"github.com/agurinov/gopl/kafka"
)

type (
	Contract struct {
		S string
		I int
		F []float64
		B bool
	}
	Consumer       = kafka.Consumer[Contract]
	ConsumerOption = kafka.ConsumerOption[Contract]
)

var (
	configmap = kafka.ConfigMap{
		kafka.BootstrapServersKey: "broker1,broker2,broker3",
		kafka.GroupIDKey:          "consumer_group_1",
		kafka.SecurityProtocolKey: "SASL_PLAINTEXT",
		kafka.SASLMechanismKey:    "PLAIN",
		kafka.SASLUsernameKey:     "sasl_user",
		kafka.SASLPasswordKey:     "sasl_password",
	}
	config = kafka.Config{
		BatchSize: 5,
		EventPosition: kafka.EventPosition{
			Topic:     "topic_1",
			Partition: 5,
			Offset:    100,
		},
	}
)

func Example() {
	var lib kafka.ComboLibrary // = segmentio.New()

	opts := []ConsumerOption{
		kafka.WithLibrary[Contract](lib),
		kafka.WithDLQ[Contract](lib),
		kafka.WithConfigMap[Contract](configmap),
		kafka.WithConsumerConfig[Contract](config),
		kafka.WithEventSerializer(kafka.JsonEventSerializer[Contract]),
		kafka.WithBackoffOptions[Contract](
			backoff.WithMaxRetries(2),
			backoff.WithStrategy(
				backoff.NewStaticStrategy(100*time.Millisecond),
			),
		),
	}

	consumer, err := kafka.NewConsumer(opts...)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	if err := consumer.Consume(ctx); err != nil {
		panic(err)
	}
}
