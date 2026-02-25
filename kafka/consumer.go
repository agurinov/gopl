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
		Ping(context.Context) error
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

func (cs consumer) Start() error {
	opts := cs.clientOptions

	opts = append(opts,
		kgo.DisableAutoCommit(),
		kgo.BlockRebalanceOnPoll(),
		kgo.OnPartitionsAssigned(cs.onAssigned),
		kgo.OnPartitionsRevoked(cs.onRevoked),
		kgo.OnPartitionsLost(cs.onRevoked),
	)

	return cs.client.Init(opts...)
}

func (cs consumer) Ping(ctx context.Context) error {
	if cs.client == nil {
		return nil
	}

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
	partitionDispatcher, err := irun.NewDispatcher[int32]()
	if err != nil {
		return nil, err
	}

	obj := consumer{
		recordMapper:        kgoRecordMapper{},
		partitionDispatcher: partitionDispatcher,
		client:              new(kgoClient),
	}

	obj, err = c.ConstructWithValidate(obj, opts...)
	if err != nil {
		return nil, err
	}

	return &obj, nil
}
