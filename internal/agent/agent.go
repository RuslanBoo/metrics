package agent

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

type Agent struct {
	serverAddr     string
	pollInterval   time.Duration
	reportInterval time.Duration
	client         *http.Client

	pollCount int64
}

func New(serverAddr string, pollInterval, reportInterval time.Duration) *Agent {
	return &Agent{
		serverAddr:     serverAddr,
		pollInterval:   pollInterval,
		reportInterval: reportInterval,
		client:         &http.Client{Timeout: 5 * time.Second},
	}
}

func (a *Agent) Run() {
	pollTicker := time.NewTicker(a.pollInterval)
	reportTicker := time.NewTicker(a.reportInterval)
	defer pollTicker.Stop()
	defer reportTicker.Stop()

	var gauges map[string]float64
	var counters map[string]int64
	gauges = make(map[string]float64)
	counters = make(map[string]int64)

	for {
		select {
		case <-pollTicker.C:
			g, c := a.CollectMetrics()
			for k, v := range g {
				gauges[k] = v
			}
			for k, v := range c {
				counters[k] = v
			}

		case <-reportTicker.C:
			a.ReportGauges(gauges)
			a.ReportCounters(counters)
		}
	}
}

func (a *Agent) CollectMetrics() (map[string]float64, map[string]int64) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	gauges := map[string]float64{
		"Alloc":         float64(mem.Alloc),
		"BuckHashSys":   float64(mem.BuckHashSys),
		"Frees":         float64(mem.Frees),
		"GCCPUFraction": mem.GCCPUFraction,
		"GCSys":         float64(mem.GCSys),
		"HeapAlloc":     float64(mem.HeapAlloc),
		"HeapIdle":      float64(mem.HeapIdle),
		"HeapInuse":     float64(mem.HeapInuse),
		"HeapObjects":   float64(mem.HeapObjects),
		"HeapReleased":  float64(mem.HeapReleased),
		"HeapSys":       float64(mem.HeapSys),
		"LastGC":        float64(mem.LastGC),
		"Lookups":       float64(mem.Lookups),
		"MCacheInuse":   float64(mem.MCacheInuse),
		"MCacheSys":     float64(mem.MCacheSys),
		"MSpanInuse":    float64(mem.MSpanInuse),
		"MSpanSys":      float64(mem.MSpanSys),
		"Mallocs":       float64(mem.Mallocs),
		"NextGC":        float64(mem.NextGC),
		"NumForcedGC":   float64(mem.NumForcedGC),
		"NumGC":         float64(mem.NumGC),
		"OtherSys":      float64(mem.OtherSys),
		"PauseTotalNs":  float64(mem.PauseTotalNs),
		"StackInuse":    float64(mem.StackInuse),
		"StackSys":      float64(mem.StackSys),
		"Sys":           float64(mem.Sys),
		"TotalAlloc":    float64(mem.TotalAlloc),
		"RandomValue":   rand.Float64(),
	}

	a.pollCount++

	counters := map[string]int64{"PollCount": a.pollCount}

	return gauges, counters
}

func (a *Agent) GetPollCount() int64 {
	return a.pollCount
}

func (a *Agent) ReportGauges(gauges map[string]float64) {
	for name, value := range gauges {
		url := fmt.Sprintf("%s/update/gauge/%s/%s", a.serverAddr, name, strconv.FormatFloat(value, 'f', -1, 64))
		req, _ := http.NewRequest(http.MethodPost, url, nil)
		req.Header.Set("Content-Type", "text/plain")
		resp, err := a.client.Do(req)
		if err == nil && resp != nil {
			resp.Body.Close()
		}
	}
}

func (a *Agent) ReportCounters(counters map[string]int64) {
	for name, value := range counters {
		url := fmt.Sprintf("%s/update/counter/%s/%d", a.serverAddr, name, value)
		req, _ := http.NewRequest(http.MethodPost, url, nil)
		req.Header.Set("Content-Type", "text/plain")
		resp, err := a.client.Do(req)
		if err == nil && resp != nil {
			resp.Body.Close()
		}
	}
}
