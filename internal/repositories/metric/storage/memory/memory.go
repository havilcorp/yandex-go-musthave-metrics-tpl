// Package memory репозиторий для работы с метриками в оперативной памяти
package memory

import (
	"context"
	"fmt"

	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/domain"
)

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

// NewMemStorage инициализация хранилища в памяти
func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauge:   make(map[string]float64, 0),
		counter: map[string]int64{},
	}
}

// AddGauge добавление метрики
func (store *MemStorage) AddGauge(ctx context.Context, key string, gauge float64) error {
	store.gauge[key] = gauge
	return nil
}

// AddCounter добавление метрики
func (store *MemStorage) AddCounter(ctx context.Context, key string, counter int64) error {
	if val, ok := store.counter[key]; ok {
		store.counter[key] = val + counter
	} else {
		store.counter[key] = counter
	}
	return nil
}

// AddGaugeBulk добавление метрики массивом
func (store *MemStorage) AddGaugeBulk(ctx context.Context, list []domain.Gauge) error {
	for _, model := range list {
		if err := store.AddGauge(ctx, model.Key, model.Value); err != nil {
			return fmt.Errorf("addGaugeBulk => %w", err)
		}
	}
	return nil
}

// AddCounterBulk добавление метрики массивом
func (store *MemStorage) AddCounterBulk(ctx context.Context, list []domain.Counter) error {
	for _, model := range list {
		if err := store.AddCounter(ctx, model.Key, model.Value); err != nil {
			return fmt.Errorf("addCounterBulk => %w", err)
		}
	}
	return nil
}

// GetGauge получение значения метрики
func (store *MemStorage) GetGauge(ctx context.Context, key string) (float64, error) {
	val, ok := store.gauge[key]
	if !ok {
		return 0, domain.ErrValueNotFound
	}
	return val, nil
}

// GetCounter получение значения метрики
func (store *MemStorage) GetCounter(ctx context.Context, key string) (int64, error) {
	val, ok := store.counter[key]
	if !ok {
		return 0, domain.ErrValueNotFound
	}
	return val, nil
}

// GetAllGauge получение всех значений метрики
func (store *MemStorage) GetAllGauge(ctx context.Context) (map[string]float64, error) {
	return store.gauge, nil
}

// GetAllCounters получение всех значений метрики
func (store *MemStorage) GetAllCounters(ctx context.Context) (map[string]int64, error) {
	return store.counter, nil
}
