package memory

import (
	"context"
	"fmt"

	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/models"
)

type MemStorage struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

func (store *MemStorage) Init() error {
	return nil
}

func (store *MemStorage) Close() {

}

func (store *MemStorage) AddGauge(ctx context.Context, key string, gauge float64) error {
	store.Gauge[key] = gauge
	return nil
}

func (store *MemStorage) AddCounter(ctx context.Context, key string, counter int64) error {
	if val, ok := store.Counter[key]; ok {
		store.Counter[key] = val + counter
	} else {
		store.Counter[key] = counter
	}
	return nil
}

func (store *MemStorage) AddGaugeBulk(ctx context.Context, list []models.GaugeModel) error {
	for _, model := range list {
		if err := store.AddGauge(ctx, model.Key, model.Value); err != nil {
			return fmt.Errorf("addGaugeBulk => %w", err)
		}
	}
	return nil
}

func (store *MemStorage) AddCounterBulk(ctx context.Context, list []models.CounterModel) error {
	for _, model := range list {
		if err := store.AddCounter(ctx, model.Key, model.Value); err != nil {
			return fmt.Errorf("addCounterBulk => %w", err)
		}
	}
	return nil
}

func (store *MemStorage) GetCounter(ctx context.Context, key string) (int64, bool) {
	val, ok := store.Counter[key]
	return val, ok
}

func (store *MemStorage) GetGauge(ctx context.Context, key string) (float64, bool) {
	val, ok := store.Gauge[key]
	return val, ok
}

func (store *MemStorage) GetAllCounters(ctx context.Context) map[string]int64 {
	return store.Counter
}

func (store *MemStorage) GetAllGauge(ctx context.Context) map[string]float64 {
	return store.Gauge
}

func (store *MemStorage) SaveToFile(ctx context.Context) error {
	return nil
}

func (store *MemStorage) Ping() error {
	return nil
}
