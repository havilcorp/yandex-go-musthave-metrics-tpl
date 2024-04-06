package file

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"time"

	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/domain"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/config"
)

type FileStorage struct {
	Conf    *config.Config
	Gauge   map[string]float64
	Counter map[string]int64
}

func NewFileStorage(conf *config.Config) (*FileStorage, error) {
	ctx := context.Background()
	fileStorage := FileStorage{
		Conf:    conf,
		Gauge:   map[string]float64{},
		Counter: map[string]int64{},
	}
	if conf.IsRestore {
		var err error
		for _, sec := range []int{1, 3, 5} {
			err = fileStorage.LoadFromFile(ctx)
			if errors.Is(err, fs.ErrNotExist) {
				return &fileStorage, nil
			}
			if errors.Is(err, fs.ErrClosed) {
				time.Sleep(time.Duration(sec) * time.Second)
			} else {
				break
			}
		}
		if err != nil {
			return nil, fmt.Errorf("init => %w", err)
		}
	}
	return &fileStorage, nil
}

func (store *FileStorage) AddGauge(ctx context.Context, key string, gauge float64) error {
	store.Gauge[key] = gauge
	if store.Conf.StoreInterval == 0 {
		if err := store.SaveToFile(ctx); err != nil {
			return fmt.Errorf("addGauge => %w", err)
		}
	}
	return nil
}

func (store *FileStorage) AddCounter(ctx context.Context, key string, counter int64) error {
	if val, ok := store.Counter[key]; ok {
		store.Counter[key] = val + counter
	} else {
		store.Counter[key] = counter
	}
	if store.Conf.StoreInterval == 0 {
		if err := store.SaveToFile(ctx); err != nil {
			return fmt.Errorf("addCounter => %w", err)
		}
	}
	return nil
}

func (store *FileStorage) AddGaugeBulk(ctx context.Context, list []domain.Gauge) error {
	for _, model := range list {
		if err := store.AddGauge(ctx, model.Key, model.Value); err != nil {
			return fmt.Errorf("addGaugeBulk => %w", err)
		}
	}
	return nil
}

func (store *FileStorage) AddCounterBulk(ctx context.Context, list []domain.Counter) error {
	for _, model := range list {
		if err := store.AddCounter(ctx, model.Key, model.Value); err != nil {
			return fmt.Errorf("addCounterBulk => %w", err)
		}
	}
	return nil
}

func (store *FileStorage) GetCounter(ctx context.Context, key string) (int64, error) {
	val, ok := store.Counter[key]
	if !ok {
		return 0, domain.ErrValueNotFound
	}
	return val, nil
}

func (store *FileStorage) GetGauge(ctx context.Context, key string) (float64, error) {
	val, ok := store.Gauge[key]
	if !ok {
		return 0, domain.ErrValueNotFound
	}
	return val, nil
}

func (store *FileStorage) GetAllCounters(ctx context.Context) (map[string]int64, error) {
	return store.Counter, nil
}

func (store *FileStorage) GetAllGauge(ctx context.Context) (map[string]float64, error) {
	return store.Gauge, nil
}

func (store *FileStorage) SaveToFile(ctx context.Context) error {
	listGauges, err := store.GetAllGauge(ctx)
	if err != nil {
		return fmt.Errorf("GetAllGauge %v", err)
	}
	listCounters, err := store.GetAllCounters(ctx)
	if err != nil {
		return fmt.Errorf("GetAllCounters %v", err)
	}
	modelFileStorage := domain.Metric{
		Gauge:   listGauges,
		Counter: listCounters,
	}
	allDataJSON, err := json.Marshal(modelFileStorage)
	if err != nil {
		return fmt.Errorf("saveToFile => %w", err)
	}
	if err := os.WriteFile(store.Conf.FileStoragePath, allDataJSON, 0o666); err != nil {
		return fmt.Errorf("saveToFile => %w", err)
	}
	return nil
}

func (store *FileStorage) LoadFromFile(ctx context.Context) error {
	file, err := os.ReadFile(store.Conf.FileStoragePath)
	if err != nil {
		return fmt.Errorf("loadFromFile => %w", err)
	}
	if len(file) < 3 {
		file = []byte("{}")
	}
	fileStorage := &FileStorage{}
	if err := json.Unmarshal(file, fileStorage); err != nil {
		return fmt.Errorf("loadFromFile => %w", err)
	}
	for key, value := range fileStorage.Counter {
		if err := store.AddCounter(ctx, key, value); err != nil {
			return fmt.Errorf("loadFromFile => %w", err)
		}
	}
	for key, value := range fileStorage.Gauge {
		if err := store.AddGauge(ctx, key, value); err != nil {
			return fmt.Errorf("loadFromFile => %w", err)
		}
	}
	return nil
}
