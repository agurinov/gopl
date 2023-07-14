//go:build test_unit

package kafka_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/agurinov/gopl/backoff"
	"github.com/agurinov/gopl/kafka"
	"github.com/agurinov/gopl/kafka/mock"
	pl_testing "github.com/agurinov/gopl/testing"
)

func TestConsumer_Validate(t *testing.T) {
	pl_testing.Init(t)

	cases := map[string]struct {
		inputConsumerOptions []ConsumerOption
		pl_testing.TestCase
	}{
		"vcase00: nil library": {
			inputConsumerOptions: []ConsumerOption{
				kafka.WithLibrary[Contract](nil),
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: kafka.ErrEmptyConsumerLibrary,
			},
		},
		"vcase01: nil event serializer": {
			inputConsumerOptions: []ConsumerOption{
				kafka.WithLibrary[Contract](new(mock.ConsumerLibrary)),
				kafka.WithEventSerializer[Contract](nil),
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: kafka.ErrEmptySerializer,
			},
		},
		"vcase02: invalid cfg": {
			inputConsumerOptions: []ConsumerOption{
				kafka.WithLibrary[Contract](new(mock.ConsumerLibrary)),
				kafka.WithEventSerializer(kafka.JsonEventSerializer[Contract]),
				kafka.WithConsumerConfig[Contract](kafka.Config{BatchSize: 0}),
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: kafka.ErrInvalidConfig,
			},
		},
		"vcase03: invalid configmap": {
			inputConsumerOptions: []ConsumerOption{
				kafka.WithLibrary[Contract](new(mock.ConsumerLibrary)),
				kafka.WithEventSerializer(kafka.JsonEventSerializer[Contract]),
				kafka.WithConsumerConfig[Contract](config),
				kafka.WithConfigMap[Contract](nil),
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: kafka.ErrInvalidConfigMap,
			},
		},
		"vcase04: nil onebyone handler": {
			inputConsumerOptions: []ConsumerOption{
				kafka.WithLibrary[Contract](new(mock.ConsumerLibrary)),
				kafka.WithEventSerializer(kafka.JsonEventSerializer[Contract]),
				kafka.WithConsumerConfig[Contract](config),
				kafka.WithConfigMap[Contract](configmap),
				kafka.WithEventHandler[Contract](nil),
				kafka.WithEventBatchHandler[Contract](new(mock.EventBatchHandler[Contract])),
				kafka.WithEventHandleStrategy[Contract](kafka.EventHandleOneByOne),
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: kafka.ErrInvalidEventHandler,
			},
		},
		"vcase05: nil batch handler": {
			inputConsumerOptions: []ConsumerOption{
				kafka.WithLibrary[Contract](new(mock.ConsumerLibrary)),
				kafka.WithEventSerializer(kafka.JsonEventSerializer[Contract]),
				kafka.WithConsumerConfig[Contract](config),
				kafka.WithConfigMap[Contract](configmap),
				kafka.WithEventHandler[Contract](new(mock.EventHandler[Contract])),
				kafka.WithEventBatchHandler[Contract](nil),
				kafka.WithEventHandleStrategy[Contract](kafka.EventHandleBatch),
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: kafka.ErrInvalidEventHandler,
			},
		},
		"vcase06: valid consumer": {
			inputConsumerOptions: []ConsumerOption{
				kafka.WithLibrary[Contract](new(mock.ConsumerLibrary)),
				kafka.WithEventSerializer(kafka.JsonEventSerializer[Contract]),
				kafka.WithConsumerConfig[Contract](config),
				kafka.WithConfigMap[Contract](configmap),
				kafka.WithEventHandler[Contract](new(mock.EventHandler[Contract])),
			},
		},
		"vcase07: unsupported handle strategy": {
			inputConsumerOptions: []ConsumerOption{
				kafka.WithLibrary[Contract](new(mock.ConsumerLibrary)),
				kafka.WithEventSerializer(kafka.JsonEventSerializer[Contract]),
				kafka.WithConsumerConfig[Contract](config),
				kafka.WithConfigMap[Contract](configmap),
				kafka.WithEventHandleStrategy[Contract](kafka.EventHandleStrategy(100)),
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: kafka.ErrInvalidEventHandler,
			},
		},
	}

	for name, tc := range cases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			consumer, err := kafka.NewConsumer[Contract](tc.inputConsumerOptions...)
			tc.CheckError(t, err)
			require.NotNil(t, consumer)
		})
	}
}

func TestConsumer_Consume(t *testing.T) {
	pl_testing.Init(t)

	cases := map[string]struct {
		inputContextTimeout  time.Duration
		inputConsumerOptions []ConsumerOption
		mocks                func(mocks)
		pl_testing.TestCase
	}{
		"case00: library init err": {
			mocks: func(m mocks) {
				gomock.InOrder(
					m.library.EXPECT().Init(IsContext(), configmap, config).Times(1).Return(errLibraryInit),
				)
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: errLibraryInit,
			},
		},
		"case01: dlq init err": {
			mocks: func(m mocks) {
				gomock.InOrder(
					m.library.EXPECT().Init(IsContext(), configmap, config).Times(1).Return(nil),
					m.dlq.EXPECT().Init(IsContext(), configmap, config).Times(1).Return(errDLQInit),
					m.library.EXPECT().Close().Times(1).Return(nil),
				)
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: errDLQInit,
			},
		},
		"case02: dlq; batch; consume err overrides close err": {
			mocks: func(m mocks) {
				gomock.InOrder(
					m.library.EXPECT().Init(IsContext(), configmap, config).Times(1).Return(nil),
					m.dlq.EXPECT().Init(IsContext(), configmap, config).Times(1).Return(nil),
					m.library.EXPECT().ConsumeBatch(IsContext(), uint(5)).Times(1).Return(
						nil,
						kafka.EmptyEventPosition,
						errConsumePermanent,
					),
					m.dlq.EXPECT().Close().Times(1).Return(errDLQClose),
					m.library.EXPECT().Close().Times(1).Return(errLibraryClose),
				)
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: errConsumePermanent,
			},
		},
		"case03: dlq; batch; consume retryable err; backoff limit": {
			mocks: func(m mocks) {
				gomock.InOrder(
					m.library.EXPECT().Init(IsContext(), configmap, config).Times(1).Return(nil),
					m.dlq.EXPECT().Init(IsContext(), configmap, config).Times(1).Return(nil),
					m.library.EXPECT().ConsumeBatch(IsContext(), uint(5)).Times(3).Return(
						nil,
						kafka.EmptyEventPosition,
						errConsumeRetryable,
					),
					m.dlq.EXPECT().Close().Times(1).Return(errDLQClose),
					m.library.EXPECT().Close().Times(1).Return(errLibraryClose),
				)
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: errConsumeRetryable,
			},
		},
		"case04: dlq; batch; consumed no events; backoff limit": {
			mocks: func(m mocks) {
				gomock.InOrder(
					m.library.EXPECT().Init(IsContext(), configmap, config).Times(1).Return(nil),
					m.dlq.EXPECT().Init(IsContext(), configmap, config).Times(1).Return(nil),
					m.library.EXPECT().ConsumeBatch(IsContext(), uint(5)).Times(3).Return(
						nil,
						kafka.EmptyEventPosition,
						nil,
					),
					m.dlq.EXPECT().Close().Times(1).Return(errDLQClose),
					m.library.EXPECT().Close().Times(1).Return(errLibraryClose),
				)
			},
			TestCase: pl_testing.TestCase{
				MustFail: true,
				MustFailIsErr: backoff.RetryLimitError{
					BackoffName: "kafka-consumer",
					MaxRetries:  2,
				},
			},
		},
		"case05: dlq; batch; consumed unexpected event position": {
			mocks: func(m mocks) {
				gomock.InOrder(
					m.library.EXPECT().Init(IsContext(), configmap, config).Times(1).Return(nil),
					m.dlq.EXPECT().Init(IsContext(), configmap, config).Times(1).Return(nil),
					m.library.EXPECT().ConsumeBatch(IsContext(), uint(5)).Times(1).Return(
						[][]byte{event5Json, event3Json},
						MutateEventPositionTopic("unexpected_topic"),
						nil,
					),
					m.dlq.EXPECT().Close().Times(1).Return(errDLQClose),
					m.library.EXPECT().Close().Times(1).Return(errLibraryClose),
				)
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: kafka.ErrInvalidEventPosition,
			},
		},
		"case06: no dlq; onebyone; serialization err; handle ok": {
			inputConsumerOptions: []ConsumerOption{
				kafka.WithDLQ[Contract](nil),
				kafka.WithEventHandleStrategy[Contract](kafka.EventHandleOneByOne),
			},
			mocks: func(m mocks) {
				gomock.InOrder(
					m.library.EXPECT().Init(IsContext(), configmap, config).Times(1).Return(nil),
					m.library.EXPECT().ConsumeBatch(IsContext(), uint(5)).Times(1).Return(
						[][]byte{event1Json, brokenEvent1Json, event2Json},
						eventPosition,
						nil,
					),
					m.eventHandler.EXPECT().Handle(IsContext(), event1).Times(1).Return(nil),
					m.library.EXPECT().Commit(
						IsContext(),
						MutateEventPositionOffset(12345),
					).Times(1).Return(nil),
					m.library.EXPECT().Close().Times(1).Return(errLibraryClose),
				)
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailAsErr: NewSerializationError(),
			},
		},
		"case07: no dlq; onebyone; serialization err; handle err": {
			inputConsumerOptions: []ConsumerOption{
				kafka.WithDLQ[Contract](nil),
				kafka.WithEventHandleStrategy[Contract](kafka.EventHandleOneByOne),
			},
			mocks: func(m mocks) {
				gomock.InOrder(
					m.library.EXPECT().Init(IsContext(), configmap, config).Times(1).Return(nil),
					m.library.EXPECT().ConsumeBatch(IsContext(), uint(5)).Times(1).Return(
						[][]byte{event1Json, brokenEvent1Json, event2Json},
						eventPosition,
						nil,
					),
					m.eventHandler.EXPECT().Handle(IsContext(), event1).Times(1).Return(errHandleEvent),
					m.library.EXPECT().Close().Times(1).Return(errLibraryClose),
				)
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailAsErr: NewSerializationError(),
			},
		},
		"case08: no dlq; onebyone; serialization err; handle err; nothing to commit": {
			inputConsumerOptions: []ConsumerOption{
				kafka.WithDLQ[Contract](nil),
				kafka.WithEventHandleStrategy[Contract](kafka.EventHandleOneByOne),
			},
			mocks: func(m mocks) {
				gomock.InOrder(
					m.library.EXPECT().Init(IsContext(), configmap, config).Times(1).Return(nil),
					m.library.EXPECT().ConsumeBatch(IsContext(), uint(5)).Times(1).Return(
						[][]byte{event3Json, event5Json, brokenEvent3Json, brokenEvent1Json},
						eventPosition,
						nil,
					),
					m.eventHandler.EXPECT().Handle(IsContext(), event3).Times(1).Return(errHandleEvent),
					m.library.EXPECT().Close().Times(1).Return(errLibraryClose),
				)
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailAsErr: NewSerializationError(),
			},
		},
		"case09: no dlq; onebyone; serialization err; handle err; least success committed": {
			inputConsumerOptions: []ConsumerOption{
				kafka.WithDLQ[Contract](nil),
				kafka.WithEventHandleStrategy[Contract](kafka.EventHandleOneByOne),
			},
			mocks: func(m mocks) {
				gomock.InOrder(
					m.library.EXPECT().Init(IsContext(), configmap, config).Times(1).Return(nil),
					m.library.EXPECT().ConsumeBatch(IsContext(), uint(5)).Times(1).Return(
						[][]byte{event4Json, event2Json, brokenEvent2Json},
						eventPosition,
						nil,
					),
					m.eventHandler.EXPECT().Handle(IsContext(), event4).Times(1).Return(nil),
					m.eventHandler.EXPECT().Handle(IsContext(), event2).Times(1).Return(errHandleEvent),
					m.library.EXPECT().Commit(
						IsContext(),
						MutateEventPositionOffset(12345),
					).Times(1).Return(nil),
					m.library.EXPECT().Close().Times(1).Return(errLibraryClose),
				)
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailAsErr: NewSerializationError(),
			},
		},
		"case10: dlq; batch; serialization err; handle ok": {
			inputConsumerOptions: []ConsumerOption{
				kafka.WithEventHandleStrategy[Contract](kafka.EventHandleBatch),
			},
			mocks: func(m mocks) {
				gomock.InOrder(
					m.library.EXPECT().Init(IsContext(), configmap, config).Times(1).Return(nil),
					m.dlq.EXPECT().Init(IsContext(), configmap, config).Times(1).Return(nil),
					m.library.EXPECT().ConsumeBatch(IsContext(), uint(5)).Times(1).Return(
						[][]byte{event1Json, brokenEvent1Json, event5Json},
						eventPosition,
						nil,
					),
					m.eventBatchHandler.EXPECT().Handle(
						IsContext(),
						[]Contract{event1, event5},
					).Times(1).Return(nil),
					m.dlq.EXPECT().ProduceBatch(IsContext(),
						brokenEvent1Json,
					).Times(1).Return(nil),
					m.library.EXPECT().Commit(
						IsContext(),
						MutateEventPositionOffset(12347),
					).Times(1).Return(nil),
					m.dlq.EXPECT().Close().Times(1).Return(errDLQClose),
					m.library.EXPECT().Close().Times(1).Return(errLibraryClose),
				)
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailAsErr: NewSerializationError(),
			},
		},
		"case11: dlq; batch; serialization err; handle err": {
			inputConsumerOptions: []ConsumerOption{
				kafka.WithEventHandleStrategy[Contract](kafka.EventHandleBatch),
			},
			mocks: func(m mocks) {
				gomock.InOrder(
					m.library.EXPECT().Init(IsContext(), configmap, config).Times(1).Return(nil),
					m.dlq.EXPECT().Init(IsContext(), configmap, config).Times(1).Return(nil),
					m.library.EXPECT().ConsumeBatch(IsContext(), uint(5)).Times(1).Return(
						[][]byte{event1Json, brokenEvent1Json, event5Json},
						eventPosition,
						nil,
					),
					m.eventBatchHandler.EXPECT().Handle(
						IsContext(),
						[]Contract{event1, event5},
					).Times(1).Return(errHandleEvent),
					m.dlq.EXPECT().Close().Times(1).Return(errDLQClose),
					m.library.EXPECT().Close().Times(1).Return(errLibraryClose),
				)
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailAsErr: NewSerializationError(),
			},
		},
		"case12: dlq; onebyone; serialization err; handle err": {
			inputConsumerOptions: []ConsumerOption{
				kafka.WithEventHandleStrategy[Contract](kafka.EventHandleOneByOne),
			},
			mocks: func(m mocks) {
				gomock.InOrder(
					m.library.EXPECT().Init(IsContext(), configmap, config).Times(1).Return(nil),
					m.dlq.EXPECT().Init(IsContext(), configmap, config).Times(1).Return(nil),
					m.library.EXPECT().ConsumeBatch(IsContext(), uint(5)).Times(1).Return(
						[][]byte{event1Json, brokenEvent1Json, event5Json},
						eventPosition,
						nil,
					),
					m.eventHandler.EXPECT().Handle(IsContext(), event1).Times(1).Return(errHandleEvent),
					m.eventHandler.EXPECT().Handle(IsContext(), event5).Times(1).Return(errHandleEvent),
					m.dlq.EXPECT().ProduceBatch(IsContext(),
						brokenEvent1Json, event1Json, event5Json,
					).Times(1).Return(nil),
					m.library.EXPECT().Commit(
						IsContext(),
						MutateEventPositionOffset(12347),
					).Times(1).Return(nil),
					m.dlq.EXPECT().Close().Times(1).Return(errDLQClose),
					m.library.EXPECT().Close().Times(1).Return(errLibraryClose),
				)
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailAsErr: NewSerializationError(),
			},
		},
		"case13: dlq; onebyone; handle err; produce retryable err; backoff limit": {
			inputConsumerOptions: []ConsumerOption{
				kafka.WithEventHandleStrategy[Contract](kafka.EventHandleOneByOne),
			},
			mocks: func(m mocks) {
				gomock.InOrder(
					m.library.EXPECT().Init(IsContext(), configmap, config).Times(1).Return(nil),
					m.dlq.EXPECT().Init(IsContext(), configmap, config).Times(1).Return(nil),
					m.library.EXPECT().ConsumeBatch(IsContext(), uint(5)).Times(1).Return(
						[][]byte{event3Json, event1Json, event5Json},
						eventPosition,
						nil,
					),
					m.eventHandler.EXPECT().Handle(IsContext(), event3).Times(1).Return(errHandleEvent),
					m.eventHandler.EXPECT().Handle(IsContext(), event1).Times(1).Return(errHandleEvent),
					m.eventHandler.EXPECT().Handle(IsContext(), event5).Times(1).Return(errHandleEvent),
					m.dlq.EXPECT().ProduceBatch(IsContext(),
						event3Json, event1Json, event5Json,
					).Times(3).Return(errDLQProduceRetryable),
					m.dlq.EXPECT().Close().Times(1).Return(errDLQClose),
					m.library.EXPECT().Close().Times(1).Return(errLibraryClose),
				)
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: errDLQProduceRetryable,
			},
		},
		"case13: dlq; onebyone; handle err; commit retryable err; backoff limit": {
			inputConsumerOptions: []ConsumerOption{
				kafka.WithEventHandleStrategy[Contract](kafka.EventHandleOneByOne),
			},
			mocks: func(m mocks) {
				gomock.InOrder(
					m.library.EXPECT().Init(IsContext(), configmap, config).Times(1).Return(nil),
					m.dlq.EXPECT().Init(IsContext(), configmap, config).Times(1).Return(nil),
					m.library.EXPECT().ConsumeBatch(IsContext(), uint(5)).Times(1).Return(
						[][]byte{event3Json, event1Json, event5Json},
						eventPosition,
						nil,
					),
					m.eventHandler.EXPECT().Handle(IsContext(), event3).Times(1).Return(errHandleEvent),
					m.eventHandler.EXPECT().Handle(IsContext(), event1).Times(1).Return(errHandleEvent),
					m.eventHandler.EXPECT().Handle(IsContext(), event5).Times(1).Return(errHandleEvent),
					m.dlq.EXPECT().ProduceBatch(IsContext(),
						event3Json, event1Json, event5Json,
					).Times(1).Return(nil),
					m.library.EXPECT().Commit(
						IsContext(),
						MutateEventPositionOffset(12347),
					).Times(3).Return(errLibraryCommitRetryable),
					m.dlq.EXPECT().Close().Times(1).Return(errDLQClose),
					m.library.EXPECT().Close().Times(1).Return(errLibraryClose),
				)
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: errLibraryCommitRetryable,
			},
		},
		"case14: no dlq; batch; finished by context": {
			inputContextTimeout: 20 * time.Millisecond,
			inputConsumerOptions: []ConsumerOption{
				kafka.WithDLQ[Contract](nil),
				kafka.WithMaxIterations[Contract](0),
			},
			mocks: func(m mocks) {
				m.library.EXPECT().Init(IsContext(), configmap, config).Times(1).Return(nil)
				m.library.EXPECT().ConsumeBatch(IsContext(), uint(5)).AnyTimes().Return(
					[][]byte{event1Json, event2Json},
					eventPosition,
					nil,
				)
				m.eventBatchHandler.EXPECT().Handle(
					IsContext(),
					[]Contract{event1, event2},
				).AnyTimes().Return(nil)
				m.library.EXPECT().Commit(
					IsContext(),
					MutateEventPositionOffset(12346),
				).AnyTimes().Return(nil)
				m.library.EXPECT().Close().Times(1).Return(errLibraryClose)
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: context.DeadlineExceeded,
			},
		},
		"case15: no dlq; onebyone; finished by max iterations": {
			inputConsumerOptions: []ConsumerOption{
				kafka.WithDLQ[Contract](nil),
				kafka.WithEventHandleStrategy[Contract](kafka.EventHandleOneByOne),
				kafka.WithMaxIterations[Contract](5),
			},
			mocks: func(m mocks) {
				m.library.EXPECT().Init(IsContext(), configmap, config).Times(1).Return(nil)
				m.library.EXPECT().ConsumeBatch(IsContext(), uint(5)).Times(5).Return(
					[][]byte{event1Json, event2Json, event5Json, event4Json},
					eventPosition,
					nil,
				)
				m.eventHandler.EXPECT().Handle(IsContext(), event1).Times(5).Return(nil)
				m.eventHandler.EXPECT().Handle(IsContext(), event2).Times(5).Return(nil)
				m.eventHandler.EXPECT().Handle(IsContext(), event5).Times(5).Return(nil)
				m.eventHandler.EXPECT().Handle(IsContext(), event4).Times(5).Return(nil)
				m.library.EXPECT().Commit(
					IsContext(),
					MutateEventPositionOffset(12348),
				).Times(5).Return(nil)
				m.library.EXPECT().Close().Times(1).Return(errLibraryClose)
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: errLibraryClose,
			},
		},

		// Complex case with dlq support
		// iteration 1: consumed=2     serialized=0/2 handled=0/0 dlq=2     committed=2
		// iteration 2: consumed=2     serialized=1/2 handled=0/1 dlq=0,0,2 committed=2
		// iteration 3: consumed=4     serialized=3/4 handled=2/3 dlq=2     committed=4
		// iteration 4: consumed=4     serialized=3/4 handled=3/3 dlq=1     committed=4
		// iteration 5: consumed=0,0,5 serialized=5/5 handled=0/5 dlq=5     committed=0,0,5 (check backoff reset)
		// iteration 6: consumed=3     serialized=3/3 handled=2/3 dlq=1     committed=3
		// iteration 7: consumed=1     serialized=1/1 handled=1/1 dlq=0     committed=1
		"case16: dlq; onebyone; complex": {
			inputConsumerOptions: []ConsumerOption{
				kafka.WithEventHandleStrategy[Contract](kafka.EventHandleOneByOne),
				kafka.WithMaxIterations[Contract](7),
			},
			mocks: func(m mocks) {
				gomock.InOrder(
					m.library.EXPECT().Init(IsContext(), configmap, config).Times(1).Return(nil),
					m.dlq.EXPECT().Init(IsContext(), configmap, config).Times(1).Return(nil),

					// iteration 1
					m.library.EXPECT().ConsumeBatch(IsContext(), uint(5)).Times(1).Return(
						[][]byte{brokenEvent1Json, brokenEvent2Json},
						eventPosition,
						nil,
					),
					m.dlq.EXPECT().ProduceBatch(IsContext(),
						brokenEvent1Json, brokenEvent2Json,
					).Times(1).Return(nil),
					m.library.EXPECT().Commit(
						IsContext(),
						MutateEventPositionOffset(12346),
					).Times(1).Return(nil),

					// iteration 2
					m.library.EXPECT().ConsumeBatch(IsContext(), uint(5)).Times(1).Return(
						[][]byte{brokenEvent2Json, event3Json},
						eventPosition,
						nil,
					),
					m.eventHandler.EXPECT().Handle(IsContext(), event3).Times(1).Return(errHandleEvent),
					m.dlq.EXPECT().ProduceBatch(IsContext(),
						brokenEvent2Json, event3Json,
					).Times(2).Return(errDLQProduceRetryable),
					m.dlq.EXPECT().ProduceBatch(IsContext(),
						brokenEvent2Json, event3Json,
					).Times(1).Return(nil),
					m.library.EXPECT().Commit(
						IsContext(),
						MutateEventPositionOffset(12346),
					).Times(1).Return(nil),

					// iteration 3
					m.library.EXPECT().ConsumeBatch(IsContext(), uint(5)).Times(1).Return(
						[][]byte{event2Json, event4Json, brokenEvent5Json, event3Json},
						eventPosition,
						nil,
					),
					m.eventHandler.EXPECT().Handle(IsContext(), event2).Times(1).Return(nil),
					m.eventHandler.EXPECT().Handle(IsContext(), event4).Times(1).Return(errHandleEvent),
					m.eventHandler.EXPECT().Handle(IsContext(), event3).Times(1).Return(nil),
					m.dlq.EXPECT().ProduceBatch(IsContext(),
						brokenEvent5Json, event4Json,
					).Times(1).Return(nil),
					m.library.EXPECT().Commit(
						IsContext(),
						MutateEventPositionOffset(12348),
					).Times(1).Return(nil),

					// iteration 4
					m.library.EXPECT().ConsumeBatch(IsContext(), uint(5)).Times(1).Return(
						[][]byte{event3Json, brokenEvent4Json, event1Json, event5Json},
						eventPosition,
						nil,
					),
					m.eventHandler.EXPECT().Handle(IsContext(), event3).Times(1).Return(nil),
					m.eventHandler.EXPECT().Handle(IsContext(), event1).Times(1).Return(nil),
					m.eventHandler.EXPECT().Handle(IsContext(), event5).Times(1).Return(nil),
					m.dlq.EXPECT().ProduceBatch(IsContext(),
						brokenEvent4Json,
					).Times(1).Return(nil),
					m.library.EXPECT().Commit(
						IsContext(),
						MutateEventPositionOffset(12348),
					).Times(1).Return(nil),

					// iteration 5
					m.library.EXPECT().ConsumeBatch(IsContext(), uint(5)).Times(2).Return(
						nil,
						kafka.EmptyEventPosition,
						errConsumeRetryable,
					),
					m.library.EXPECT().ConsumeBatch(IsContext(), uint(5)).Times(1).Return(
						[][]byte{event3Json, event5Json, event1Json, event4Json, event2Json},
						eventPosition,
						nil,
					),
					m.eventHandler.EXPECT().Handle(IsContext(), event3).Times(1).Return(errHandleEvent),
					m.eventHandler.EXPECT().Handle(IsContext(), event5).Times(1).Return(errHandleEvent),
					m.eventHandler.EXPECT().Handle(IsContext(), event1).Times(1).Return(errHandleEvent),
					m.eventHandler.EXPECT().Handle(IsContext(), event4).Times(1).Return(errHandleEvent),
					m.eventHandler.EXPECT().Handle(IsContext(), event2).Times(1).Return(errHandleEvent),
					m.dlq.EXPECT().ProduceBatch(IsContext(),
						event3Json, event5Json, event1Json, event4Json, event2Json,
					).Times(1).Return(nil),
					m.library.EXPECT().Commit(
						IsContext(),
						MutateEventPositionOffset(12349),
					).Times(2).Return(errLibraryCommitRetryable),
					m.library.EXPECT().Commit(
						IsContext(),
						MutateEventPositionOffset(12349),
					).Times(1).Return(nil),

					// iteration 6
					m.library.EXPECT().ConsumeBatch(IsContext(), uint(5)).Times(1).Return(
						[][]byte{event5Json, event1Json, event2Json},
						eventPosition,
						nil,
					),
					m.eventHandler.EXPECT().Handle(IsContext(), event5).Times(1).Return(errHandleEvent),
					m.eventHandler.EXPECT().Handle(IsContext(), event1).Times(1).Return(nil),
					m.eventHandler.EXPECT().Handle(IsContext(), event2).Times(1).Return(nil),
					m.dlq.EXPECT().ProduceBatch(IsContext(),
						event5Json,
					).Times(1).Return(nil),
					m.library.EXPECT().Commit(
						IsContext(),
						MutateEventPositionOffset(12347),
					).Times(1).Return(nil),

					// iteration 7
					m.library.EXPECT().ConsumeBatch(IsContext(), uint(5)).Times(1).Return(
						[][]byte{event4Json},
						eventPosition,
						nil,
					),
					m.eventHandler.EXPECT().Handle(IsContext(), event4).Times(1).Return(nil),
					m.library.EXPECT().Commit(
						IsContext(),
						MutateEventPositionOffset(12345),
					).Times(1).Return(nil),

					m.dlq.EXPECT().Close().Times(1).Return(nil),
					m.library.EXPECT().Close().Times(1).Return(nil),
				)
			},
			TestCase: pl_testing.TestCase{
				MustFail: false,
			},
		},
		"case17: dlq; onebyone; serialization err; handle ok; produce err": {
			inputConsumerOptions: []ConsumerOption{
				kafka.WithEventHandleStrategy[Contract](kafka.EventHandleOneByOne),
			},
			mocks: func(m mocks) {
				gomock.InOrder(
					m.library.EXPECT().Init(IsContext(), configmap, config).Times(1).Return(nil),
					m.dlq.EXPECT().Init(IsContext(), configmap, config).Times(1).Return(nil),
					m.library.EXPECT().ConsumeBatch(IsContext(), uint(5)).Times(1).Return(
						[][]byte{event1Json, brokenEvent1Json},
						eventPosition,
						nil,
					),
					m.eventHandler.EXPECT().Handle(IsContext(), event1).Times(1).Return(nil),
					m.dlq.EXPECT().ProduceBatch(IsContext(),
						brokenEvent1Json,
					).Times(1).Return(errDLQProducePermanent),
					m.dlq.EXPECT().Close().Times(1).Return(errDLQClose),
					m.library.EXPECT().Close().Times(1).Return(errLibraryClose),
				)
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailAsErr: NewSerializationError(),
			},
		},
		"case18: dlq; batch; commit err": {
			inputConsumerOptions: []ConsumerOption{
				kafka.WithEventHandleStrategy[Contract](kafka.EventHandleBatch),
			},
			mocks: func(m mocks) {
				gomock.InOrder(
					m.library.EXPECT().Init(IsContext(), configmap, config).Times(1).Return(nil),
					m.dlq.EXPECT().Init(IsContext(), configmap, config).Times(1).Return(nil),
					m.library.EXPECT().ConsumeBatch(IsContext(), uint(5)).Times(1).Return(
						[][]byte{event1Json, event5Json, event4Json},
						eventPosition,
						nil,
					),
					m.eventBatchHandler.EXPECT().Handle(
						IsContext(),
						[]Contract{event1, event5, event4},
					).Times(1).Return(nil),

					m.library.EXPECT().Commit(
						IsContext(),
						MutateEventPositionOffset(12347),
					).Times(1).Return(errLibraryCommitPermanent),
					m.dlq.EXPECT().Close().Times(1).Return(errDLQClose),
					m.library.EXPECT().Close().Times(1).Return(errLibraryClose),
				)
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: errLibraryCommitPermanent,
			},
		},
	}

	for name, tc := range cases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			m := NewMocks(ctrl)
			if tc.mocks != nil {
				tc.mocks(m)
			}

			opts := []ConsumerOption{
				kafka.WithLibrary[Contract](m.library),
				kafka.WithDLQ[Contract](m.dlq),
				kafka.WithEventSerializer(kafka.JsonEventSerializer[Contract]),
				kafka.WithConfigMap[Contract](configmap),
				kafka.WithConsumerConfig[Contract](config),
				kafka.WithEventHandler[Contract](m.eventHandler),
				kafka.WithEventBatchHandler[Contract](m.eventBatchHandler),
				kafka.WithMaxIterations[Contract](1),
				kafka.WithBackoffOptions[Contract](
					backoff.WithMaxRetries(2),
					backoff.WithStrategy(
						backoff.NewStaticStrategy(time.Millisecond),
					),
				),
			}
			opts = append(opts, tc.inputConsumerOptions...)

			consumer, err := kafka.NewConsumer[Contract](opts...)
			require.NoError(t, err)
			require.NotNil(t, consumer)

			ctx := context.Background()
			if tc.inputContextTimeout != 0 {
				var cancel context.CancelFunc

				ctx, cancel = context.WithTimeout(ctx, tc.inputContextTimeout)
				t.Cleanup(cancel)
			}

			tc.CheckError(t, consumer.Consume(ctx))
		})
	}
}
