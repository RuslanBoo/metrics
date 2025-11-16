package config

import "time"

type AgentConfig struct {
	Addr           string
	PollInterval   time.Duration
	ReportInterval time.Duration
}
