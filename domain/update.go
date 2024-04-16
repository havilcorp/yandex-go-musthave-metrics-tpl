package domain

const (
	TypeMetricsGauge   = "gauge"
	TypeMetricsCounter = "counter"
)

type MetricRequest struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}
