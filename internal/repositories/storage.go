// Package repositories репозиторий для сохранения данных метрик
package repositories

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/domain"
)

type Storage struct {
	path string
}

func NewStorage(path string) *Storage {
	return &Storage{
		path: path,
	}
}

func (store *Storage) SaveToFile(metric domain.Metric) error {
	allDataJSON, err := json.Marshal(metric)
	if err != nil {
		return fmt.Errorf("saveToFile => %w", err)
	}
	if err := os.WriteFile(store.path, allDataJSON, 0o666); err != nil {
		return fmt.Errorf("saveToFile => %w", err)
	}
	return nil
}
