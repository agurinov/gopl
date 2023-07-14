package kafka_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"reflect"

	"github.com/golang/mock/gomock"

	"github.com/agurinov/gopl/kafka"
	"github.com/agurinov/gopl/kafka/mock"
)

type mocks struct {
	library           *mock.ConsumerLibrary
	dlq               *mock.ProducerLibrary
	eventHandler      *mock.EventHandler[Contract]
	eventBatchHandler *mock.EventBatchHandler[Contract]
}

func NewMocks(ctrl *gomock.Controller) mocks {
	return mocks{
		library:           mock.NewConsumerLibrary(ctrl),
		dlq:               mock.NewProducerLibrary(ctrl),
		eventHandler:      mock.NewEventHandler[Contract](ctrl),
		eventBatchHandler: mock.NewEventBatchHandler[Contract](ctrl),
	}
}

var (
	errLibraryInit            = errors.New("can't init library")
	errLibraryClose           = errors.New("can't close library")
	errLibraryCommitRetryable = fmt.Errorf("can't commit event position: %w", context.DeadlineExceeded)
	errLibraryCommitPermanent = errors.New("can't commit event position; permanent err")
	errConsumeRetryable       = fmt.Errorf("can't consume: %w", io.EOF)
	errConsumePermanent       = errors.New("can't consume; permanent err")
	errHandleEvent            = errors.New("can't handle event")
	errDLQInit                = errors.New("can't init dlq")
	errDLQClose               = errors.New("can't close dlq")
	errDLQProduceRetryable    = fmt.Errorf("can't produce broken events to dlq: %w", &net.DNSError{IsTemporary: true})
	errDLQProducePermanent    = errors.New("can't produce broken events to dlq; permanent err")
)

var (
	event1     = Contract{S: "event1", I: 10, F: []float64{0.111, 1.222, 2.333}, B: true}
	event1Json = []byte(`{"s":"event1","i":10,"f":[0.111,1.222,2.333],"b":true,"unknown_field":"a"}`)

	event2     = Contract{S: "event2", I: 20, B: false}
	event2Json = []byte(`{"s":"event2","i":20,"unknown_field":"a"}`)

	event3     = Contract{S: "event3", I: 30, F: []float64{1.2345}, B: true}
	event3Json = []byte(`{"s":"event3","i":30,"f":[1.2345],"b":true,"unknown_field":"a"}`)

	event4     = Contract{S: "event4", I: 40, F: []float64{0.01, 0.99}}
	event4Json = []byte(`{"s":"event4","i":40,"f":[0.01,0.99],"unknown_field":"a"}`)

	event5     = Contract{S: "event5", I: 50, F: []float64{}, B: true}
	event5Json = []byte(`{"s":"event5","i":50,"f":[],"b":true,"unknown_field":"a"}`)

	brokenEvent1Json = []byte("foo")
	brokenEvent2Json = []byte("bar")
	brokenEvent3Json = []byte("baz")
	brokenEvent4Json = []byte("lol")
	brokenEvent5Json = []byte("kek")
)

var eventPosition = kafka.EventPosition{
	Topic:     "topic_1",
	Partition: 5,
	Offset:    12345,
}

func IsContext() gomock.Matcher {
	return gomock.AssignableToTypeOf(
		reflect.TypeOf((*context.Context)(nil)).Elem(),
	)
}

func MutateEventPositionTopic(topic string) kafka.EventPosition {
	return kafka.EventPosition{
		Topic:     topic,
		Partition: eventPosition.Partition,
		Offset:    eventPosition.Offset,
	}
}

func MutateEventPositionOffset(offset int64) kafka.EventPosition {
	return kafka.EventPosition{
		Topic:     eventPosition.Topic,
		Partition: eventPosition.Partition,
		Offset:    offset,
	}
}

func MutateConfigBatchSize(s uint) kafka.Config {
	return kafka.Config{
		BatchSize:     s,
		EventPosition: config.EventPosition,
	}
}

func NewSerializationError() interface{} {
	serializationErr := new(json.SyntaxError)

	return &serializationErr
}
