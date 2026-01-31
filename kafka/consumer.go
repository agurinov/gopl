package kafka

import (
	"context"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/agurinov/gopl/backoff"
	"github.com/agurinov/gopl/graceful"
	irun "github.com/agurinov/gopl/internal/run"
	c "github.com/agurinov/gopl/patterns/creational"
	"github.com/agurinov/gopl/run"
)

type (
	Consumer interface {
		Start() error
		Close(context.Context) error
	}
	consumer struct {
		metrics             consumerMetrics
		recordDiscarder     RecordDiscarder
		recordMapper        RecordMapper
		partitionDispatcher irun.Dispatcher[int32]
		logger              *zap.Logger
		client              consumerClient
		clientOptions       []kgo.Opt
		handler             Handler
		handlerBatch        HandlerBatch
		topic               string
		backoffFabric       backoff.Fabric
		maxPollDuration     time.Duration
		maxPollRecords      int
	}
	ConsumerOption c.Option[consumer]
)

func (cs consumer) Close(ctx context.Context) error {
	closeFn := run.SimpleFn(func() {
		cs.client.Close()
	})

	return graceful.Close(closeFn)(ctx)
}

func (cs *consumer) Start() error {
	client, err := newKgoClient(cs.clientOptions...)
	if err != nil {
		return err
	}

	cs.client = client

	return nil
}

func (cs consumer) Ping(ctx context.Context) error {
	err := cs.client.Ping(ctx)

	logLvl := zapcore.InfoLevel
	if err != nil {
		logLvl = zapcore.ErrorLevel
	}

	cs.logger.Log(
		logLvl,
		"consumer ping",
		zap.String("topic", cs.topic),
		zap.Int32s("assigned.partitions", cs.partitionDispatcher.Running()),
		zap.Error(err),
	)

	return err
}

func NewConsumer(
	opts ...ConsumerOption,
) (
	Consumer,
	error,
) {
	obj := consumer{
		recordMapper: kgoRecordMapper{},
	}

	opts = append(opts,
		WithConsumerClientOptions(
			kgo.DisableAutoCommit(),
			kgo.BlockRebalanceOnPoll(),
			kgo.OnPartitionsAssigned(obj.onAssigned),
			kgo.OnPartitionsRevoked(obj.onRevoked),
			kgo.OnPartitionsLost(obj.onRevoked),
		),
	)

	obj, err := c.ConstructWithValidate(obj, opts...)
	if err != nil {
		return nil, err
	}

	partitionDispatcher, err := irun.NewDispatcher[int32]()
	if err != nil {
		return nil, err
	}

	obj.partitionDispatcher = partitionDispatcher

	return &obj, nil
}
