package web

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

func TestMetricsWithAttributes(t *testing.T) {
	ctx := context.Background()
	mReader := metric.NewManualReader()
	defer func() {
		assert.NoError(t, mReader.Shutdown(ctx))
	}()
	testProvider := metric.NewMeterProvider(metric.WithReader(mReader))
	testMeter := testProvider.Meter("test-meter")

	mm := NewMetricsMiddleware(testMeter, "")
	testServer := httptest.NewServer(mm.Handle(testingHandler))
	defer testServer.Close()

	_, err := http.Get(testServer.URL)
	assert.NoError(t, err)
	testRequestTelemetryAttr := attribute.NewSet(
		semconv.HTTPRequestMethodKey.String("GET"),
		semconv.HTTPResponseStatusCode(http.StatusOK),
		semconv.HTTPRouteKey.String("/"),
	)

	var rmData metricdata.ResourceMetrics
	err = mReader.Collect(ctx, &rmData)
	assert.NoError(t, err)

	for _, mData := range rmData.ScopeMetrics {
		expectedMetricData := map[string]metricdata.Aggregation{
			"http.server.request.duration": metricdata.Histogram[int64]{
				Temporality: metricdata.CumulativeTemporality,
				DataPoints: []metricdata.HistogramDataPoint[int64]{
					{Count: 1, Attributes: testRequestTelemetryAttr},
				},
			},
			string(semconv.HTTPRequestBodySizeKey): metricdata.Sum[int64]{
				Temporality: metricdata.CumulativeTemporality,
				IsMonotonic: true,
				DataPoints: []metricdata.DataPoint[int64]{
					{Attributes: testRequestTelemetryAttr},
				},
			},
			string(semconv.HTTPResponseBodySizeKey): metricdata.Sum[int64]{
				Temporality: metricdata.CumulativeTemporality,
				IsMonotonic: true,
				DataPoints: []metricdata.DataPoint[int64]{
					{Attributes: testRequestTelemetryAttr},
				},
			},
		}
		for _, data := range mData.Metrics {
			if expectedMetricData[data.Name] == nil {
				continue // discard test for metrics that are not listed in `expectedMetricData`
			}
			metricdatatest.AssertAggregationsEqual(t,
				expectedMetricData[data.Name],
				data.Data,
				metricdatatest.IgnoreTimestamp(),
				metricdatatest.IgnoreValue(),
				metricdatatest.IgnoreExemplars(),
			)
			delete(expectedMetricData, data.Name)
		}
		assert.Equal(t, len(expectedMetricData), 0, "Some expected metrics are not set: %+v", reflect.ValueOf(expectedMetricData).MapKeys())
	}
}

func TestMetricsPatternOverride(t *testing.T) {
	ctx := context.Background()
	mReader := metric.NewManualReader()
	defer func() {
		assert.NoError(t, mReader.Shutdown(ctx))
	}()
	testProvider := metric.NewMeterProvider(metric.WithReader(mReader))
	testMeter := testProvider.Meter("test-meter")

	pattern := "/character/{name}"
	mm := NewMetricsMiddleware(testMeter, pattern)
	testServer := httptest.NewServer(mm.Handle(testingHandler))
	defer testServer.Close()

	_, err := http.Get(testServer.URL)
	assert.NoError(t, err)
	testRequestTelemetryAttr := []attribute.KeyValue{
		semconv.HTTPRouteKey.String(pattern),
	}

	var rmData metricdata.ResourceMetrics
	err = mReader.Collect(ctx, &rmData)
	assert.NoError(t, err)

	for _, mData := range rmData.ScopeMetrics {
		metricdatatest.AssertHasAttributes(t, mData, testRequestTelemetryAttr...)
	}
}
