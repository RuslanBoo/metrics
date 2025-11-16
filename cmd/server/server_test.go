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

func setupHandler() (*handler.Handler, *memory.MemStorage) {
	metricRepository := memory.New()
	metricService := service.New(metricRepository)
	return handler.New(metricService), metricRepository
}

func postToHandler(h *handler.Handler, url string, t *testing.T) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, url, nil)
	w := httptest.NewRecorder()
	h.Update(w, req)
	return w
}

func TestUpdateHandler_StatusCodes(t *testing.T) {
	metricHandler, metricRepository := setupHandler()

	sendAndCheck := func(url string) {
		w := postToHandler(metricHandler, url, t)
		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}

	sendAndCheck("/update/gauge/testGauge/99.99")
	val, ok := metricRepository.GetGauge("testGauge")
	assert.True(t, ok, "Gauge should exist")
	assert.Equal(t, 99.99, float64(val), "Gauge value should match") // приведение к float64

	sendAndCheck("/update/counter/testCounter/3")
	sendAndCheck("/update/counter/testCounter/2")
	cVal, ok := metricRepository.GetCounter("testCounter")
	assert.True(t, ok, "Counter should exist")
	assert.Equal(t, int64(5), int64(cVal), "Counter value should sum correctly") // приведение к int64
}

func TestMetricsStoredCorrectly(t *testing.T) {
	metricHandler, metricRepository := setupHandler()

	sendAndCheck := func(url string) {
		w := postToHandler(metricHandler, url, t)
		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}

	sendAndCheck("/update/gauge/testGauge/99.99")
	val, ok := metricRepository.GetGauge("testGauge")
	assert.True(t, ok, "Gauge should exist")
	assert.Equal(t, 99.99, float64(val), "Gauge value should match") // приведение к float64

	sendAndCheck("/update/counter/testCounter/3")
	sendAndCheck("/update/counter/testCounter/2")
	cVal, ok := metricRepository.GetCounter("testCounter")
	assert.True(t, ok, "Counter should exist")
	assert.Equal(t, int64(5), int64(cVal), "Counter value should sum correctly") // приведение к int64
}

func TestMultipleMetrics(t *testing.T) {
	metricHandler, metricRepository := setupHandler()

	metrics := map[string]float64{
		"Alloc":       12345,
		"HeapAlloc":   67890,
		"RandomValue": 0.42,
	}

	for name, val := range metrics {
		url := fmt.Sprintf("/update/gauge/%s/%f", name, val)
		w := postToHandler(metricHandler, url, t)
		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		v, ok := metricRepository.GetGauge(name)
		assert.True(t, ok, fmt.Sprintf("Gauge %s should exist", name))
		assert.Equal(t, val, float64(v), "Gauge value should match") // приведение к float64
	}
}
