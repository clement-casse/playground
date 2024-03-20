package web

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"

	"github.com/stretchr/testify/assert"
)

var (
	testingHandler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, "test ok")
	})

	panicReason     = "something weird"
	panicingHandler = http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		panic(panicReason)
	})
)

func TestRecoveryMiddleware(t *testing.T) {
	for _, tt := range []struct {
		name         string
		handlerFunc  http.HandlerFunc
		expectStatus int
	}{
		{
			name:         "a normal handler should not be affected",
			handlerFunc:  testingHandler,
			expectStatus: http.StatusOK,
		}, {
			name:         "a panic in the handler should be caught",
			handlerFunc:  panicingHandler,
			expectStatus: http.StatusInternalServerError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			rm := NewRecoveryMiddleware(slog.Default(), nil)
			testServer := httptest.NewServer(rm.Handle(tt.handlerFunc))
			defer testServer.Close()

			res, err := http.Get(testServer.URL)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectStatus, res.StatusCode, "unexpected status")
		})
	}
}

func TestRecoveryMiddlewareLogsPanicReason(t *testing.T) {
	var recorder bytes.Buffer
	strLogger := slog.New(slog.NewTextHandler(&recorder, nil))

	rm := NewRecoveryMiddleware(strLogger, nil)
	testServer := httptest.NewServer(rm.Handle(panicingHandler))

	_, err := http.Get(testServer.URL)
	assert.NoError(t, err)
	testServer.Close()

	logLines := recorder.String()
	assert.Contains(t, logLines, panicReason, "recovery middleware does not print the inner panic reason")
}

func TestRecoveryMiddlewareIncrementsCounter(t *testing.T) {
	var recorder bytes.Buffer
	strLogger := slog.New(slog.NewTextHandler(&recorder, nil))

	ctx := context.Background()
	mReader := metric.NewManualReader()
	defer func() {
		assert.NoError(t, mReader.Shutdown(ctx))
	}()
	testProvider := metric.NewMeterProvider(metric.WithReader(mReader))
	testMeter := testProvider.Meter("test-meter")

	rm := NewRecoveryMiddleware(strLogger, testMeter)
	testServer := httptest.NewServer(rm.Handle(panicingHandler))
	defer testServer.Close()

	var err error
	_, err = http.Get(testServer.URL)
	assert.NoError(t, err)
	_, err = http.Get(testServer.URL)
	assert.NoError(t, err)

	var rmData metricdata.ResourceMetrics
	err = mReader.Collect(ctx, &rmData)
	assert.NoError(t, err)

	for _, mData := range rmData.ScopeMetrics {
		expectedMetricData := map[string]metricdata.Aggregation{
			"application.panics.recovered": metricdata.Sum[int64]{
				Temporality: metricdata.CumulativeTemporality,
				IsMonotonic: true,
				DataPoints: []metricdata.DataPoint[int64]{
					{Value: 2},
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
				metricdatatest.IgnoreExemplars(),
			)
			delete(expectedMetricData, data.Name)
		}
		assert.Equal(t, len(expectedMetricData), 0, "Some expected metrics are not set: %+v", reflect.ValueOf(expectedMetricData).MapKeys())
	}
}
