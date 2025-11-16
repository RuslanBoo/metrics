package service

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/RuslanBoo/metrics/internal/config"
	"github.com/RuslanBoo/metrics/internal/model"
	"github.com/RuslanBoo/metrics/internal/repository"
)

var (
	ErrUnknownMetricType = errors.New("unknown metric type")
	ErrInvalidValue      = errors.New("invalid metric value")
)

type MetricService struct {
	repository repository.Storage
}

func New(repo repository.Storage) *MetricService {
	return &MetricService{repository: repo}
}

func (s *MetricService) UpdateMetric(metricType, name, valueStr string) error {
	switch metricType {
	case string(models.GaugeType):
		val, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			return ErrInvalidValue
		}
		s.repository.SetGauge(name, models.Gauge(val))
	case string(models.CounterType):
		delta, err := strconv.ParseInt(valueStr, 10, 64)
		if err != nil {
			return ErrInvalidValue
		}
		s.repository.AddCounter(name, models.Counter(delta))
	default:
		return ErrUnknownMetricType
	}
	return nil
}

func (s *MetricService) GetMetric(metricType, name string) (interface{}, bool) {
	switch metricType {
	case string(models.GaugeType):
		val, ok := s.repository.GetGauge(name)
		return val, ok
	case string(models.CounterType):
		val, ok := s.repository.GetCounter(name)
		return val, ok
	default:
		return nil, false
	}
}

func (s *MetricService) GetAllMetrics() (map[string]models.Gauge, map[string]models.Counter) {
	return s.repository.GetAllGauges(), s.repository.GetAllCounters()
}

func ParseAgentFlags() config.AgentConfig {
	var addr string
	var pollSec int
	var reportSec int

	flag.StringVar(&addr, "a", "http://localhost:8080", "Server address")
	flag.IntVar(&pollSec, "p", 2, "Poll interval in seconds")
	flag.IntVar(&reportSec, "r", 10, "Report interval in seconds")

	flag.Parse()

	if flag.NArg() > 0 {
		fmt.Fprintf(os.Stderr, "Unknown flag(s): %v\n", flag.Args())
		os.Exit(1)
	}

	return config.AgentConfig{
		Addr:           addr,
		PollInterval:   time.Duration(pollSec) * time.Second,
		ReportInterval: time.Duration(reportSec) * time.Second,
	}
}

func ParseServerFlags() config.ServerConfig {
	var addr string

	flag.StringVar(&addr, "a", ":8080", "Server address")
	flag.Parse()

	if flag.NArg() > 0 {
		fmt.Fprintf(os.Stderr, "Unknown flag(s): %v\n", flag.Args())
		os.Exit(1)
	}

	return config.ServerConfig{
		Addr: addr,
	}
}
