package segmentio

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	segmentio "github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"

	pl_errors "github.com/agurinov/gopl/errors"
	"github.com/agurinov/gopl/kafka"
)

type impl struct {
	reader *segmentio.Reader
	writer *segmentio.Writer

	consumerMode bool
	producerMode bool

	mu sync.Mutex
}

func (i *impl) Init(ctx context.Context, cm kafka.ConfigMap, cfg kafka.Config) (err error) {
	if !i.consumerMode && !i.producerMode {
		err = fmt.Errorf("any from consumer or producer mode must be set")

		return
	}

	var (
		dialer = &segmentio.Dialer{
			DualStack: true,
		}
		protocol = cm[kafka.SecurityProtocolKey]
	)

	var (
		needSSL = protocol == kafka.SASLSSLSecurityProtocol ||
			protocol == kafka.SSLSecurityProtocol
		needSASL = protocol == kafka.SASLPlaintextSecurityProtocol ||
			protocol == kafka.SASLSSLSecurityProtocol
	)

	if needSSL {
		certPool := x509.NewCertPool()
		if !certPool.AppendCertsFromPEM(
			[]byte(cm[kafka.SSLCAPEMKey]),
		) {
			err = fmt.Errorf("can't parse PEM")

			return
		}

		dialer.TLS = &tls.Config{
			MinVersion:         tls.VersionTLS12,
			RootCAs:            certPool,
			InsecureSkipVerify: true,
		}
	}

	if needSASL {
		switch mechanism := cm[kafka.SASLMechanismKey]; mechanism {
		case kafka.PlainSASLMechanism:
			dialer.SASLMechanism = plain.Mechanism{
				Username: cm[kafka.SASLUsernameKey],
				Password: cm[kafka.SASLPasswordKey],
			}
		default:
			err = fmt.Errorf("%w: %s", kafka.ErrUnknownSASLMechanism, mechanism)

			return
		}
	}

	var (
		brokers = strings.Split(cm[kafka.BootstrapServersKey], ",")
		dialErr error
	)

	for _, broker := range brokers {
		var (
			addr = segmentio.TCP(broker)
			conn *segmentio.Conn
		)

		if conn, dialErr = dialer.DialContext(ctx, addr.Network(), addr.String()); dialErr == nil {
			defer conn.Close()

			break
		}
	}

	if dialErr != nil {
		err = dialErr

		return
	}

	// segmentio vendor panics on creation
	defer func() {
		if r := recover(); r != nil {
			if asErr, ok := r.(error); ok {
				err = asErr
			} else {
				panic(r)
			}
		}
	}()

	i.mu.Lock()
	defer i.mu.Unlock()

	if i.consumerMode {
		i.reader = segmentio.NewReader(segmentio.ReaderConfig{
			Brokers:        brokers,
			GroupID:        cm[kafka.GroupIDKey],
			Topic:          cfg.Topic,
			StartOffset:    segmentio.FirstOffset,
			QueueCapacity:  int(cfg.BatchSize),
			IsolationLevel: segmentio.ReadCommitted,
			MinBytes:       10e3, // 10KB
			MaxBytes:       10e6, // 10MB
			MaxWait:        300 * time.Millisecond,
			Dialer:         dialer,
			CommitInterval: 0, // Commits must be synchronous
			GroupBalancers: []segmentio.GroupBalancer{
				SingleAssignmentGroupBalancer{
					Topic:     cfg.Topic,
					Partition: int(cfg.Partition),
				},
			},
			Logger: log.New(os.Stdout, "kafka reader: ", 0),
		})

		/*
			o.consumerGroup = ????
			group, err := kafka.NewConsumerGroup(kafka.ConsumerGroupConfig{
					ID:      "my-group",
						Brokers: []string{"kafka:9092"},
							Topics:  []string{"my-topic"},
						})
		*/
	}

	if i.producerMode {
		i.writer = segmentio.NewWriter(segmentio.WriterConfig{
			Brokers:       brokers,
			Topic:         cfg.Topic,
			Dialer:        dialer,
			QueueCapacity: int(cfg.BatchSize),
			Balancer:      &segmentio.LeastBytes{},
			RequiredAcks:  int(segmentio.RequireAll),
			Async:         false,
			// Logger:        log.New(os.Stdout, "kafka writer: ", 0),
		})
	}

	return
}

func (i *impl) ConsumeBatch(
	ctx context.Context,
	batchSize uint,
) (
	[][]byte,
	kafka.EventPosition,
	error,
) {
	if !i.consumerMode {
		return nil, kafka.EmptyEventPosition, fmt.Errorf("segmentio is not configured as consumer")
	}

	var (
		events   = make([][]byte, 0, batchSize)
		position kafka.EventPosition
	)

	if offset := i.reader.Offset(); offset != -1 {
		return nil, kafka.EmptyEventPosition, fmt.Errorf(
			"%w: offset is not backed by consumer group",
			kafka.ErrInvalidEventPosition,
		)
	}

loop:
	for j := 0; j < int(batchSize); j++ {
		message, err := i.reader.FetchMessage(ctx)
		fmt.Println("MSG: ", message, err)

		switch {
		case errors.Is(err, context.DeadlineExceeded):
			break loop
		case err != nil:
			return nil, kafka.EmptyEventPosition, err
		}

		iterPosition := kafka.EventPosition{
			Topic:     message.Topic,
			Partition: int32(message.Partition),
			Offset:    message.Offset,
		}

		if j == 0 {
			position = iterPosition
		} else {
			if err := iterPosition.ValidateWith(position); err != nil {
				return nil, kafka.EmptyEventPosition, fmt.Errorf("unexpected event: %w", err)
			}
		}

		events = append(events, message.Value)
	}

	return events, position, nil
}

func (i *impl) Commit(ctx context.Context, position kafka.EventPosition) error {
	if !i.consumerMode {
		return fmt.Errorf("segmentio is not configured as consumer")
	}

	var (
		expectedTopic = i.reader.Config().Topic
		message       = segmentio.Message{
			Topic:     expectedTopic,
			Partition: int(position.Partition),
			Offset:    position.Offset,
		}
	)

	if position.Topic != "" && position.Topic != expectedTopic {
		return fmt.Errorf(
			"%w: trying to commit to an unexpected topic; %q instead of %q",
			kafka.ErrInvalidEventPosition,
			position.Topic,
			expectedTopic,
		)
	}

	if err := i.reader.CommitMessages(ctx, message); err != nil {
		return err
	}

	return nil
}

func (i *impl) ProduceBatch(ctx context.Context, events ...[]byte) error {
	if !i.producerMode {
		return fmt.Errorf("segmentio is not configured as producer")
	}

	messages := make([]segmentio.Message, 0, len(events))

	for i := range events {
		messages = append(messages, segmentio.Message{
			Value: events[i],
		})
	}

	if err := i.writer.WriteMessages(ctx, messages...); err != nil {
		return err
	}

	return nil
}

func (i *impl) Close() (err error) {
	if i.consumerMode {
		if consumerCloseErr := i.reader.Close(); consumerCloseErr != nil {
			err = pl_errors.Or(err, consumerCloseErr)
		}
	}

	if i.producerMode {
		if producerCloseErr := i.writer.Close(); producerCloseErr != nil {
			err = pl_errors.Or(err, producerCloseErr)
		}
	}

	return nil
}

func New() kafka.ComboLibrary {
	return &impl{
		consumerMode: true,
		producerMode: true,
	}
}

func NewProducer() kafka.ProducerLibrary {
	return &impl{
		producerMode: true,
	}
}

func NewConsumer() kafka.ConsumerLibrary {
	return &impl{
		consumerMode: true,
	}
}
