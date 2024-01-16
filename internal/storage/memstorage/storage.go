package memstorage

import (
	"encoding/json"
	"os"
)

type Repositories interface {
	AddCounter(key string, counter int64) error
	AddGauge(key string, gauge float64) error
	GetCounter(key string) (int64, bool)
	GetGauge(key string) (float64, bool)
	GetAllCounters() map[string]int64
	GetAllGauge() map[string]float64
	SetWfiteFileName(filename string)
	SetSyncWrite(isSync bool)
	SaveToFile() error
	LoadFromFile() error
}

type MemStorage struct {
	Gauge     map[string]float64
	Counter   map[string]int64
	syncWrite bool
	fileName  string
}

func NewMemStorage(syncWrite bool) *MemStorage {
	ms := MemStorage{Gauge: map[string]float64{}, Counter: map[string]int64{}}
	ms.SetSyncWrite(syncWrite)
	return &ms
}

func (ms *MemStorage) SetWfiteFileName(filename string) {
	ms.fileName = filename
}

func (ms *MemStorage) SetSyncWrite(isSync bool) {
	ms.syncWrite = isSync
}

func (ms *MemStorage) SaveToFile() error {
	file, err := os.OpenFile(ms.fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	allDataJSON, err := json.Marshal(ms)
	if err != nil {
		return err
	}
	if _, err := file.Write([]byte(allDataJSON)); err != nil {
		return err
	}
	return nil
}

func (ms *MemStorage) LoadFromFile() error {
	file, err := os.ReadFile(ms.fileName)
	if err != nil {
		return err
	}
	memStorage := &MemStorage{}
	if err := json.Unmarshal(file, memStorage); err != nil {
		return err
	}
	for key, value := range memStorage.Counter {
		if err := ms.AddCounter(key, value); err != nil {
			return err
		}
	}
	for key, value := range memStorage.Gauge {
		if err := ms.AddGauge(key, value); err != nil {
			return err
		}
	}
	return nil
}

func (ms *MemStorage) AddGauge(key string, gauge float64) error {
	ms.Gauge[key] = gauge
	if ms.syncWrite {
		if err := ms.SaveToFile(); err != nil {
			return err
		}
	}
	return nil
}

func (ms *MemStorage) AddCounter(key string, counter int64) error {
	if val, ok := ms.Counter[key]; ok {
		ms.Counter[key] = val + counter
	} else {
		ms.Counter[key] = counter
	}
	if ms.syncWrite {
		if err := ms.SaveToFile(); err != nil {
			return err
		}
	}
	return nil
}

func (ms *MemStorage) GetCounter(key string) (int64, bool) {
	val, ok := ms.Counter[key]
	return val, ok
}

func (ms *MemStorage) GetGauge(key string) (float64, bool) {
	val, ok := ms.Gauge[key]
	return val, ok
}

func (ms *MemStorage) GetAllCounters() map[string]int64 {
	return ms.Counter
}

func (ms *MemStorage) GetAllGauge() map[string]float64 {
	return ms.Gauge
}
