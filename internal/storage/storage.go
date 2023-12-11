package storage

import (
	"errors"
)

type Repositories interface {
	AddCounter(key string, counter int64) error
	AddGauge(key string, gauge float64) error
	GetCounter(key string) (int64, error)
	GetGauge(key string) (int64, error)
}

type MemStorage struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

// type Guage string

var availableMtricTypes = []string{
	"Alloc",
	"BuckHashSys",
	"Frees",
	"GCCPUFraction",
	"GCSys",
	"HeapAlloc",
	"HeapIdle",
	"HeapInuse",
	"HeapObjects",
	"HeapReleased",
	"HeapSys",
	"LastGC",
	"Lookups",
	"MCacheInuse",
	"MSpanSys",
	"Mallocs",
	"NextGC",
	"NumForcedGC",
	"NumGC",
	"OtherSys",
	"PauseTotalNs",
	"StackInuse",
	"StackSys",
	"Sys",
	"TotalAlloc",
	"PollCount",
	"RandomValue",
}

// const (
// 	Alloc         Guage = "Alloc"
// 	BuckHashSys   Guage = "BuckHashSys"
// 	Frees         Guage = "Frees"
// 	GCCPUFraction Guage = "GCCPUFraction"
// )

// func stringInSlice(a string, list []string) bool {
// 	for _, b := range list {
// 		if b == a {
// 			return true
// 		}
// 	}
// 	return false
// }

// func (ms MemStorage) AddGauge(key string, gauge float64) error {
// 	if ok := stringInSlice(key, availableMtricTypes); ok {
// 		ms.Gauge[key] = gauge
// 		return nil
// 	} else {
// 		return errors.New("metric type is not supported")
// 	}
// }

// func (ms MemStorage) AddCounter(key string, counter int64) error {
// 	if ok := stringInSlice(key, availableMtricTypes); ok {
// 		if val, ok := ms.Counter[key]; ok {
// 			ms.Counter[key] = val + counter
// 		} else {
// 			ms.Counter[key] = counter
// 		}
// 		return nil
// 	} else {
// 		return errors.New("metric type is not supported")
// 	}
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
