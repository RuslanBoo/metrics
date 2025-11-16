package main

import (
	"github.com/RuslanBoo/metrics/internal/agent"
	"github.com/RuslanBoo/metrics/internal/service"
)

func main() {
	metricConfig := service.ParseAgentFlags()
	metricAgent := agent.New(metricConfig.Addr, metricConfig.PollInterval, metricConfig.ReportInterval)
	metricAgent.Run()
}
