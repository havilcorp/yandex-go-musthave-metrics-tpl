package storage

type Repositories interface {
	AddGauge(key string, gauge float64) error
	AddCounter(key string, counter int64) error
}

type MemStorage struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

// type Guage string

// const (
// 	Alloc         Guage = "Alloc"
// 	BuckHashSys   Guage = "BuckHashSys"
// 	Frees         Guage = "Frees"
// 	GCCPUFraction Guage = "GCCPUFraction"
// )

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
