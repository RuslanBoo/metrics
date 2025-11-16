package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RuslanBoo/metrics/internal/handler"
	"github.com/RuslanBoo/metrics/internal/repository/memory"
	"github.com/RuslanBoo/metrics/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupHandler() (*handler.Handler, *memory.MemStorage, *gin.Engine) {
	metricRepository := memory.New()
	metricService := service.New(metricRepository)
	metricHandler := handler.New(metricService)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	metricHandler.RegisterRoutes(r)

	return metricHandler, metricRepository, r
}

func postToHandler(r *gin.Engine, url string, t *testing.T) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, url, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestUpdateHandler_StatusCodes(t *testing.T) {
	_, metricRepository, router := setupHandler()

	sendAndCheck := func(url string) {
		w := postToHandler(router, url, t)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	sendAndCheck("/update/gauge/testGauge/99.99")
	val, ok := metricRepository.GetGauge("testGauge")
	assert.True(t, ok, "Gauge should exist")
	assert.Equal(t, float64(99.99), float64(val), "Gauge value should match")

	sendAndCheck("/update/counter/testCounter/3")
	sendAndCheck("/update/counter/testCounter/2")
	cVal, ok := metricRepository.GetCounter("testCounter")
	assert.True(t, ok, "Counter should exist")
	assert.Equal(t, int64(5), int64(cVal), "Counter value should sum correctly")
}
