package models

type MetricModel struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

type GaugeModel struct {
	Key   string
	Value float64
}

type CounterModel struct {
	Key   string
	Value int64
}
