package memory

import (
	"context"

	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/models"
)

type MemStorage struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

func (store *MemStorage) Init(ctx context.Context) error {
	return nil
}

func (store *MemStorage) Close() {

}

func (store *MemStorage) AddGauge(key string, gauge float64) error {
	store.Gauge[key] = gauge
	return nil
}

func (store *MemStorage) AddCounter(key string, counter int64) error {
	if val, ok := store.Counter[key]; ok {
		store.Counter[key] = val + counter
	} else {
		store.Counter[key] = counter
	}
	return nil
}

func (store *MemStorage) AddGaugeBulk(list []models.GaugeModel) error {
	for _, model := range list {
		if err := store.AddGauge(model.Key, model.Value); err != nil {
			return err
		}
	}
	return nil
}

func (store *MemStorage) AddCounterBulk(list []models.CounterModel) error {
	for _, model := range list {
		if err := store.AddCounter(model.Key, model.Value); err != nil {
			return err
		}
	}
	return nil
}

func (store *MemStorage) GetCounter(key string) (int64, bool) {
	val, ok := store.Counter[key]
	return val, ok
}

func (store *MemStorage) GetGauge(key string) (float64, bool) {
	val, ok := store.Gauge[key]
	return val, ok
}

func (store *MemStorage) GetAllCounters() map[string]int64 {
	return store.Counter
}

func (store *MemStorage) GetAllGauge() map[string]float64 {
	return store.Gauge
}

func (store *MemStorage) SaveToFile() error {
	return nil
}

func (store *MemStorage) Ping() error {
	return nil
}
