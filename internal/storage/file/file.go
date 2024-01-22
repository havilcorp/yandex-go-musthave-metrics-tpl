package file

import (
	"context"
	"encoding/json"
	"os"

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
		if err := store.LoadFromFile(); err != nil {
			return err
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
			return err
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
			return err
		}
	}
	return nil
}

func (store *FileStorage) AddGaugeBulk(list []models.GaugeModel) error {
	for _, model := range list {
		if err := store.AddGauge(model.Key, model.Value); err != nil {
			return err
		}
	}
	return nil
}

func (store *FileStorage) AddCounterBulk(list []models.CounterModel) error {
	for _, model := range list {
		if err := store.AddCounter(model.Key, model.Value); err != nil {
			return err
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
		return err
	}
	if err := os.WriteFile(store.Conf.FileStoragePath, allDataJSON, 0666); err != nil {
		return err
	}
	return nil
}

func (store *FileStorage) LoadFromFile() error {
	file, err := os.ReadFile(store.Conf.FileStoragePath)
	if err != nil {
		return err
	}
	if len(file) < 3 {
		file = []byte("{}")
	}
	fileStorage := &FileStorage{}
	if err := json.Unmarshal(file, fileStorage); err != nil {
		return err
	}
	for key, value := range fileStorage.Counter {
		if err := store.AddCounter(key, value); err != nil {
			return err
		}
	}
	for key, value := range fileStorage.Gauge {
		if err := store.AddGauge(key, value); err != nil {
			return err
		}
	}
	return nil
}

func (store *FileStorage) Ping() error {
	return nil
}
