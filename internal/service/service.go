package service

import (
	"errors"
	"strconv"

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
	return &MetricService{
		repository: repo,
	}
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
