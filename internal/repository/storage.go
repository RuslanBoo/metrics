package repository

import models "github.com/RuslanBoo/metrics/internal/model"

type Storage interface {
	SetGauge(name string, value models.Gauge)
	AddCounter(name string, delta models.Counter)
	GetGauge(name string) (models.Gauge, bool)
	GetCounter(name string) (models.Counter, bool)
	GetAllGauges() map[string]models.Gauge
	GetAllCounters() map[string]models.Counter
}
