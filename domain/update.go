// Package domain модели входных данных
package domain

const (
	TypeMetricsGauge   = "gauge"
	TypeMetricsCounter = "counter"
)

type MetricRequest struct {
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
	ID    string   `json:"id"`
	MType string   `json:"type"`
}
