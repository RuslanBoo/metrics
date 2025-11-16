package main

import (
	"time"

	"github.com/RuslanBoo/metrics/internal/agent"
)

func main() {
	metricAgent := agent.New("http://localhost:8080", 2*time.Second, 10*time.Second)
	metricAgent.Run()
}
