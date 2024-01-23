package file

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"time"

	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/config"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/models"
)

type FileStorage struct {
	Conf    *config.Config
	Gauge   map[string]float64
	Counter map[string]int64
}

func (store *FileStorage) Init(ctx context.Context) error {
	if store.Conf.IsRestore {
		var err error
		for _, sec := range []int{1, 3, 5} {
			err = store.LoadFromFile()
			if errors.Is(err, fs.ErrClosed) {
				time.Sleep(time.Duration(sec) * time.Second)
			} else {
				break
			}
		}
		if err != nil {
			return fmt.Errorf("init => %w", err)
		}
	}
	return nil
}

func (store *FileStorage) Close() {

}

func (store *FileStorage) AddGauge(key string, gauge float64) error {
	store.Gauge[key] = gauge
	if store.Conf.StoreInterval == 0 {
		if err := store.SaveToFile(); err != nil {
			return fmt.Errorf("addGauge => %w", err)
		}
	}
	return nil
}

func (store *FileStorage) AddCounter(key string, counter int64) error {
	if val, ok := store.Counter[key]; ok {
		store.Counter[key] = val + counter
	} else {
		store.Counter[key] = counter
	}
	if store.Conf.StoreInterval == 0 {
		if err := store.SaveToFile(); err != nil {
			return fmt.Errorf("addCounter => %w", err)
		}
	}
	return nil
}

func (store *FileStorage) AddGaugeBulk(list []models.GaugeModel) error {
	for _, model := range list {
		if err := store.AddGauge(model.Key, model.Value); err != nil {
			return fmt.Errorf("addGaugeBulk => %w", err)
		}
	}
	return nil
}

func (store *FileStorage) AddCounterBulk(list []models.CounterModel) error {
	for _, model := range list {
		if err := store.AddCounter(model.Key, model.Value); err != nil {
			return fmt.Errorf("addCounterBulk => %w", err)
		}
	}
	return nil
}

func (store *FileStorage) GetCounter(key string) (int64, bool) {
	val, ok := store.Counter[key]
	return val, ok
}

func (store *FileStorage) GetGauge(key string) (float64, bool) {
	val, ok := store.Gauge[key]
	return val, ok
}

func (store *FileStorage) GetAllCounters() map[string]int64 {
	return store.Counter
}

func (store *FileStorage) GetAllGauge() map[string]float64 {
	return store.Gauge
}

func (store *FileStorage) SaveToFile() error {
	modelFileStorage := models.MetricModel{
		Gauge:   store.GetAllGauge(),
		Counter: store.GetAllCounters(),
	}
	allDataJSON, err := json.Marshal(modelFileStorage)
	if err != nil {
		return fmt.Errorf("saveToFile => %w", err)
	}
	if err := os.WriteFile(store.Conf.FileStoragePath, allDataJSON, 0666); err != nil {
		return fmt.Errorf("saveToFile => %w", err)
	}
	return nil
}

func (store *FileStorage) LoadFromFile() error {
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
		if err := store.AddCounter(key, value); err != nil {
			return fmt.Errorf("loadFromFile => %w", err)
		}
	}
	for key, value := range fileStorage.Gauge {
		if err := store.AddGauge(key, value); err != nil {
			return fmt.Errorf("loadFromFile => %w", err)
		}
	}
	return nil
}

func (store *FileStorage) Ping() error {
	return nil
}
