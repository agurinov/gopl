package kafka_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"
	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/agurinov/gopl/diag/metrics"
	"github.com/agurinov/gopl/kafka"
	pl_testing "github.com/agurinov/gopl/testing"
)

func TestCooldownMiddleware(t *testing.T) {
	pl_testing.Init(t)

	const (
		futureDelta = 200 * time.Millisecond
		pastDelta   = -time.Hour
	)

	pastProcessAfter, err := time.Now().Add(pastDelta).MarshalBinary()
	require.NoError(t, err)

	futureProcessAfter, err := time.Now().Add(futureDelta).MarshalBinary()
	require.NoError(t, err)

	hist := metrics.NewHistogram(
		kafka.KafkaConsumerCooldownDurationHistogramName,
		[]string{"topic", "partition"},
		metrics.WithoutServicePrefix(),
		metrics.WithUseExisting(),
	)

	type (
		args struct {
			record *kgo.Record
		}
		results struct {
			handlerWork []byte
			sampleCount uint64
			sampleSum   float64
		}
	)

	cases := map[string]struct {
		args    args
		results results
		pl_testing.TestCase
	}{
		"case00: no cooldown header": {
			args: args{
				record: &kgo.Record{
					Topic: "case00",
					Value: []byte("case00"),
				},
			},
			results: results{
				handlerWork: []byte("case00"),
				sampleCount: 0,
			},
		},
		"case01: cooldown header in the past": {
			args: args{
				record: &kgo.Record{
					Topic: "case01",
					Value: []byte("case01"),
					Headers: []kgo.RecordHeader{
						{Key: "gopl_cooldown_process_after", Value: pastProcessAfter},
					},
				},
			},
			results: results{
				handlerWork: []byte("case01"),
				sampleCount: 0,
			},
		},
		"case02: cooldown header in the future": {
			args: args{
				record: &kgo.Record{
					Topic: "case02",
					Value: []byte("case02"),
					Headers: []kgo.RecordHeader{
						{Key: "gopl_cooldown_process_after", Value: futureProcessAfter},
					},
				},
			},
			results: results{
				handlerWork: []byte("case02"),
				sampleCount: 1,
				sampleSum:   futureDelta.Seconds(),
			},
		},
		"case03: malformed cooldown header": {
			args: args{
				record: &kgo.Record{
					Topic: "case03",
					Value: []byte("case03"),
					Headers: []kgo.RecordHeader{
						{Key: "gopl_cooldown_process_after", Value: []byte("not-a-time")},
					},
				},
			},
			results: results{
				handlerWork: []byte("case03"),
				sampleCount: 0,
			},
		},
	}

	for name := range cases {
		tc := cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			ctx := t.Context()

			tc.args.record.Topic += "_" + pl_testing.RandomHash(t)

			var buf bytes.Buffer

			next := func(_ context.Context, r kafka.Record) error {
				buf.Write(r.Value)

				return nil
			}

			handler := kafka.CooldownMiddleware()(next)

			err := handler(ctx, tc.args.record)
			tc.CheckError(t, err)

			observer, err := hist.GetMetricWith(prometheus.Labels{
				"topic":     tc.args.record.Topic,
				"partition": "0",
			})
			require.NoError(t, err)
			require.NotNil(t, observer)

			metric, ok := observer.(prometheus.Metric)
			require.True(t, ok)
			require.NotNil(t, metric)

			dto, err := metrics.DTO(metric)
			require.NoError(t, err)
			require.NotNil(t, dto)

			require.Equal(t, tc.results.handlerWork, buf.Bytes())
			require.Equal(t, tc.results.sampleCount, dto.GetHistogram().GetSampleCount())

			if tc.results.sampleCount > 0 {
				require.InDelta(t, tc.results.sampleSum, dto.GetHistogram().GetSampleSum(), 0.05)
			}
		})
	}
}
