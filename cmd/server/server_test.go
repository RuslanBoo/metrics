package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RuslanBoo/metrics/internal/handler"
	"github.com/RuslanBoo/metrics/internal/repository/memory"
	"github.com/RuslanBoo/metrics/internal/service"
	"github.com/stretchr/testify/assert"
)

func setupHandler() *handler.Handler {
	metricRepository := memory.New()
	metricService := service.New(metricRepository)
	return handler.New(metricService)
}

func postToHandler(h *handler.Handler, url string, t *testing.T) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, url, nil)
	w := httptest.NewRecorder()
	h.Update(w, req)
	return w
}

func TestUpdateHandler_StatusCodes(t *testing.T) {
	handler := setupHandler()

	tests := []struct {
		url        string
		wantStatus int
	}{
		{"/update/gauge/testGauge/12.34", http.StatusOK},
		{"/update/counter/testCounter/5", http.StatusOK},
		{"/update/unknown/test/1", http.StatusBadRequest},
		{"/update/gauge/testGauge/notanumber", http.StatusBadRequest},
	}

	for _, tt := range tests {
		w := postToHandler(handler, tt.url, t)
		assert.Equal(t, tt.wantStatus, w.Result().StatusCode, "URL: %s", tt.url)
	}
}

func TestMetricsStoredCorrectly(t *testing.T) {
	metricHandler := setupHandler()
	metricRepository := memory.New()

	sendAndCheck := func(url string) {
		w := postToHandler(metricHandler, url, t)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	}

	sendAndCheck("/update/gauge/testGauge/99.99")
	val, ok := metricRepository.GetGauge("testGauge")
	assert.True(t, ok, "Gauge should exist")
	assert.Equal(t, 99.99, val, "Gauge value should match")

	sendAndCheck("/update/counter/testCounter/3")
	sendAndCheck("/update/counter/testCounter/2")
	cVal, ok := metricRepository.GetCounter("testCounter")
	assert.True(t, ok, "Counter should exist")
	assert.Equal(t, int64(5), cVal, "Counter value should sum correctly")
}

func TestMultipleMetrics(t *testing.T) {
	metricHandler := setupHandler()
	metricRepository := memory.New()

	metrics := map[string]float64{
		"Alloc":       12345,
		"HeapAlloc":   67890,
		"RandomValue": 0.42,
	}

	for name, val := range metrics {
		url := fmt.Sprintf("/update/gauge/%s/%f", name, val)
		w := postToHandler(metricHandler, url, t)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)

		v, ok := metricRepository.GetGauge(name)
		assert.True(t, ok)
		assert.Equal(t, val, v)
	}
}
