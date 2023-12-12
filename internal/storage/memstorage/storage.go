package memstorage

import (
	"errors"
)

type Repositories interface {
	AddCounter(key string, counter int64) error
	AddGauge(key string, gauge float64) error
	GetCounter(key string) (int64, error)
	GetGauge(key string) (float64, error)
	GetAllCounters() map[string]int64
	GetAllGauge() map[string]float64
}

type MemStorage struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

// type Guage string

// var availableMtricTypes = []string{
// 	"Alloc",
// 	"BuckHashSys",
// 	"Frees",
// 	"GCCPUFraction",
// 	"GCSys",
// 	"HeapAlloc",
// 	"HeapIdle",
// 	"HeapInuse",
// 	"HeapObjects",
// 	"HeapReleased",
// 	"HeapSys",
// 	"LastGC",
// 	"Lookups",
// 	"MCacheInuse",
// 	"MSpanSys",
// 	"Mallocs",
// 	"NextGC",
// 	"NumForcedGC",
// 	"NumGC",
// 	"OtherSys",
// 	"PauseTotalNs",
// 	"StackInuse",
// 	"StackSys",
// 	"Sys",
// 	"TotalAlloc",
// 	"PollCount",
// 	"RandomValue",
// }

func (ms MemStorage) AddGauge(key string, gauge float64) error {
	ms.Gauge[key] = gauge
	return nil
}

func (ms MemStorage) AddCounter(key string, counter int64) error {
	if val, ok := ms.Counter[key]; ok {
		ms.Counter[key] = val + counter
	} else {
		ms.Counter[key] = counter
	}
	return nil
}

func (ms MemStorage) GetCounter(key string) (int64, error) {
	if val, ok := ms.Counter[key]; ok {
		return val, nil
	} else {
		return 0, errors.New("metric not found")
	}
}

func (ms MemStorage) GetGauge(key string) (float64, error) {
	if val, ok := ms.Gauge[key]; ok {
		return val, nil
	} else {
		return 0, errors.New("metric not found")
	}
}

func (ms MemStorage) GetAllCounters() map[string]int64 {
	return ms.Counter
}

func (ms MemStorage) GetAllGauge() map[string]float64 {
	return ms.Gauge
}
