// Package domain модели метрик
package domain

type Metric struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

type Gauge struct {
	Key   string
	Value float64
}

type Counter struct {
	Key   string
	Value int64
}
