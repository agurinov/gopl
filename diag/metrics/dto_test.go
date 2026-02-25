package metrics_test

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/require"

	"github.com/agurinov/gopl/diag/metrics"
	pl_testing "github.com/agurinov/gopl/testing"
)

func TestDTO(t *testing.T) {
	pl_testing.Init(t)

	var (
		labels   = []string{"a", "b"}
		okLabels = prometheus.Labels{
			"a": "a",
			"b": "b",
		}
		fakeLabels = prometheus.Labels{
			"a": "bar",
			"b": "kek",
		}
	)

	t.Run("histogram", func(t *testing.T) {
		pl_testing.Init(t)

		hist := metrics.NewHistogram(
			"test_histogram",
			labels,
			metrics.WithUseExisting(),
		)
		require.NotNil(t, hist)

		const expectedSamples = 5

		for i := range expectedSamples {
			hist.With(okLabels).Observe(float64(i))
		}

		var (
			okDTO   = getHistogramDTO(t, hist, okLabels)
			fakeDTO = getHistogramDTO(t, hist, fakeLabels)
		)

		require.Equal(t, uint64(expectedSamples), okDTO.GetSampleCount())
		require.Equal(t, uint64(0), fakeDTO.GetSampleCount())
	})

	t.Run("counter", func(t *testing.T) {
		pl_testing.Init(t)

		counter := metrics.NewCounter(
			"test_counter",
			labels,
			metrics.WithUseExisting(),
		)
		require.NotNil(t, counter)

		const expectedSamples = 5

		for range expectedSamples {
			counter.With(okLabels).Inc()
		}

		var (
			okDTO   = getCounterDTO(t, counter, okLabels)
			fakeDTO = getCounterDTO(t, counter, fakeLabels)
		)

		require.Equal(t, float64(expectedSamples), okDTO.GetValue())
		require.Equal(t, float64(0), fakeDTO.GetValue())
	})

	t.Run("gauge", func(t *testing.T) {
		pl_testing.Init(t)

		gauge := metrics.NewGauge(
			"test_gauge",
			labels,
			metrics.WithUseExisting(),
		)
		require.NotNil(t, gauge)

		const expectedSamples = 5

		for range expectedSamples {
			gauge.With(okLabels).Inc()
		}

		var (
			okDTO   = getGaugeDTO(t, gauge, okLabels)
			fakeDTO = getGaugeDTO(t, gauge, fakeLabels)
		)

		require.Equal(t, float64(expectedSamples), okDTO.GetValue())
		require.Equal(t, float64(0), fakeDTO.GetValue())
	})
}

func getHistogramDTO(
	t *testing.T,
	vec *prometheus.HistogramVec,
	labels prometheus.Labels,
) *dto.Histogram {
	t.Helper()

	observer, err := vec.GetMetricWith(labels)
	require.NoError(t, err)
	require.NotNil(t, observer)

	metric, ok := observer.(prometheus.Metric)
	require.True(t, ok)
	require.NotNil(t, metric)

	dto, err := metrics.DTO(metric)
	require.NoError(t, err)
	require.NotNil(t, dto)

	return dto.GetHistogram()
}

func getCounterDTO(
	t *testing.T,
	vec *prometheus.CounterVec,
	labels prometheus.Labels,
) *dto.Counter {
	t.Helper()

	metric, err := vec.GetMetricWith(labels)
	require.NoError(t, err)
	require.NotNil(t, metric)

	dto, err := metrics.DTO(metric)
	require.NoError(t, err)
	require.NotNil(t, dto)

	return dto.GetCounter()
}

func getGaugeDTO(
	t *testing.T,
	vec *prometheus.GaugeVec,
	labels prometheus.Labels,
) *dto.Gauge {
	t.Helper()

	metric, err := vec.GetMetricWith(labels)
	require.NoError(t, err)
	require.NotNil(t, metric)

	dto, err := metrics.DTO(metric)
	require.NoError(t, err)
	require.NotNil(t, dto)

	return dto.GetGauge()
}
