package storage

import (
	"context"

	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/models"
)

type IStorage interface {
	Init() error
	Close()
	AddGauge(ctx context.Context, key string, gauge float64) error
	AddCounter(ctx context.Context, key string, counter int64) error
	AddGaugeBulk(ctx context.Context, list []models.GaugeModel) error
	AddCounterBulk(ctx context.Context, list []models.CounterModel) error
	GetCounter(ctx context.Context, key string) (int64, bool)
	GetGauge(ctx context.Context, key string) (float64, bool)
	GetAllCounters(ctx context.Context) map[string]int64
	GetAllGauge(ctx context.Context) map[string]float64
	SaveToFile(ctx context.Context) error
	Ping() error
}
