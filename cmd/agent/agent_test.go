package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/RuslanBoo/metrics/internal/agent"
	"github.com/stretchr/testify/assert"
)

func testAgent(serverURL string) *agent.Agent {
	return agent.New(serverURL, time.Millisecond, time.Millisecond)
}

func TestReportGauges(t *testing.T) {
	called := false

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true

		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "text/plain", r.Header.Get("Content-Type"))

		match, _ := regexp.MatchString(`/update/gauge/TestValue/1\.23$`, r.URL.Path)
		assert.True(t, match, "invalid gauge URL: "+r.URL.Path)
	}))
	defer testServer.Close()

	testAgent := testAgent(testServer.URL)
	testAgent.ReportGauges(map[string]float64{
		"TestValue": 1.23,
	})

	assert.True(t, called, "server must receive request")
}

func TestReportCounters(t *testing.T) {
	called := false

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true

		assert.Equal(t, http.MethodPost, r.Method)

		match, _ := regexp.MatchString(`/update/counter/PollCount/42$`, r.URL.Path)
		assert.True(t, match, "invalid counter URL: "+r.URL.Path)
	}))
	defer testServer.Close()

	testAgent := testAgent(testServer.URL)
	testAgent.ReportCounters(map[string]int64{
		"PollCount": 42,
	})

	assert.True(t, called, "server must receive counter request")
}

func TestPollCountIncrement(t *testing.T) {
	testAgent := testAgent("http://localhost")

	before := testAgent.GetPollCount()
	testAgent.CollectMetrics()
	after := testAgent.GetPollCount()

	assert.Equal(t, before+1, after)
}

func TestRuntimeMetricsCollected(t *testing.T) {
	testAgent := testAgent("http://localhost")
	gauges, _ := testAgent.CollectMetrics()

	expectedMetrics := []string{
		"Alloc", "HeapIdle", "TotalAlloc", "RandomValue",
	}

	for _, key := range expectedMetrics {
		val, ok := gauges[key]
		assert.True(t, ok, fmt.Sprintf("%s must be present", key))

		if key != "RandomValue" {
			assert.GreaterOrEqual(t, val, float64(0))
		} else {
			assert.GreaterOrEqual(t, val, 0.0)
			assert.LessOrEqual(t, val, 1.0)
		}
	}
}

func TestReportAllMetrics(t *testing.T) {
	requests := 0

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
	}))
	defer testServer.Close()

	testAgent := testAgent(testServer.URL)

	gauges := map[string]float64{
		"A": 1.1,
		"B": 2.2,
		"C": 3.3,
	}
	counters := map[string]int64{
		"PollCount": 10,
		"Errors":    5,
	}

	testAgent.ReportGauges(gauges)
	testAgent.ReportCounters(counters)

	assert.Equal(t, 5, requests) // 3 gauges + 2 counters
}

func TestGaugeSerialization(t *testing.T) {
	called := false
	expected := strconv.FormatFloat(3.1415926, 'f', -1, 64)

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true

		assert.Contains(t, r.URL.String(), expected)
	}))
	defer testServer.Close()

	testAgent := testAgent(testServer.URL)
	testAgent.ReportGauges(map[string]float64{"PI": 3.1415926})

	assert.True(t, called)
}

func TestAgentSendEvenIfServerErrors(t *testing.T) {
	calls := 0

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		http.Error(w, "boom", 500)
	}))
	defer testServer.Close()

	testAgent := testAgent(testServer.URL)

	testAgent.ReportGauges(map[string]float64{"X": 1.0})
	testAgent.ReportCounters(map[string]int64{"Y": 1})

	assert.Equal(t, 2, calls, "requests must still be made even if server errors")
}
