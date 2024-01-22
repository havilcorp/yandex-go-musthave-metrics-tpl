package storage

import "context"

type IStorage interface {
	Init(ctx context.Context) error
	Close()
	AddGauge(key string, gauge float64) error
	AddCounter(key string, counter int64) error
	GetCounter(key string) (int64, bool)
	GetGauge(key string) (float64, bool)
	GetAllCounters() map[string]int64
	GetAllGauge() map[string]float64
	SaveToFile() error
	Ping() error
}
