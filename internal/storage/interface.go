package storage

import (
	"context"

	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/models"
)

type IStorage interface {
	Init(ctx context.Context) error
	Close()
	AddGauge(key string, gauge float64) error
	AddCounter(key string, counter int64) error
	AddGaugeBulk(list []models.GaugeModel) error
	AddCounterBulk(list []models.CounterModel) error
	GetCounter(key string) (int64, bool)
	GetGauge(key string) (float64, bool)
	GetAllCounters() map[string]int64
	GetAllGauge() map[string]float64
	SaveToFile() error
	Ping() error
}
