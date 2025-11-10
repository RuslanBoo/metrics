package memory

import (
	models "github.com/RuslanBoo/metrics/internal/model"
)

type MemStorage struct {
	gauges   map[string]models.Gauge
	counters map[string]models.Counter
}

func New() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]models.Gauge),
		counters: make(map[string]models.Counter),
	}
}

func (s *MemStorage) SetGauge(name string, value models.Gauge) {
	s.gauges[name] = value
}

func (s *MemStorage) AddCounter(name string, delta models.Counter) {
	s.counters[name] += delta
}

func (s *MemStorage) GetGauge(name string) (models.Gauge, bool) {
	v, ok := s.gauges[name]
	return v, ok
}

func (s *MemStorage) GetCounter(name string) (models.Counter, bool) {
	v, ok := s.counters[name]
	return v, ok
}
