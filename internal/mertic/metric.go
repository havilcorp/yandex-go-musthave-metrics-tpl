package mertic

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"math/rand"
	"runtime"

	"github.com/go-resty/resty/v2"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/models"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/storage/memory"
)

var memStats runtime.MemStats

func WriteMetric(ms memory.MemStorage) error {
	runtime.ReadMemStats(&memStats)
	gauges := map[string]float64{
		"Alloc":         float64(memStats.Alloc),
		"BuckHashSys":   float64(memStats.BuckHashSys),
		"Frees":         float64(memStats.Frees),
		"GCCPUFraction": float64(memStats.GCCPUFraction),
		"GCSys":         float64(memStats.GCSys),
		"HeapAlloc":     float64(memStats.HeapAlloc),
		"HeapIdle":      float64(memStats.HeapIdle),
		"HeapInuse":     float64(memStats.HeapInuse),
		"HeapObjects":   float64(memStats.HeapObjects),
		"HeapReleased":  float64(memStats.HeapReleased),
		"HeapSys":       float64(memStats.HeapSys),
		"LastGC":        float64(memStats.LastGC),
		"Lookups":       float64(memStats.Lookups),
		"MCacheInuse":   float64(memStats.MCacheInuse),
		"MCacheSys":     float64(memStats.MCacheSys),
		"MSpanInuse":    float64(memStats.MSpanInuse),
		"MSpanSys":      float64(memStats.MSpanSys),
		"Mallocs":       float64(memStats.Mallocs),
		"NextGC":        float64(memStats.NextGC),
		"NumForcedGC":   float64(memStats.NumForcedGC),
		"NumGC":         float64(memStats.NumGC),
		"OtherSys":      float64(memStats.OtherSys),
		"PauseTotalNs":  float64(memStats.PauseTotalNs),
		"StackInuse":    float64(memStats.StackInuse),
		"StackSys":      float64(memStats.StackSys),
		"Sys":           float64(memStats.Sys),
		"TotalAlloc":    float64(memStats.TotalAlloc),
		"RandomValue":   float64(rand.Intn(10)),
	}
	for key, val := range gauges {
		if err := ms.AddGauge(key, val); err != nil {
			return fmt.Errorf("writeMetric => %w", err)
		}
	}
	if err := ms.AddCounter("PollCount", int64(1)); err != nil {
		return fmt.Errorf("writeMetric => %w", err)
	}
	return nil
}

func SendMetric(address string, ms memory.MemStorage) error {
	client := resty.New()
	gauges := ms.GetAllGauge()
	counters := ms.GetAllCounters()
	metrics := make([]models.MetricsRequest, 0)
	url := fmt.Sprintf("http://%s/updates", address)
	buf := bytes.NewBuffer(nil)
	zb := gzip.NewWriter(buf)
	for key, val := range gauges {
		value := val
		metrics = append(metrics, models.MetricsRequest{ID: key, MType: "gauge", Value: &value})
	}
	for key, val := range counters {
		value := val
		metrics = append(metrics, models.MetricsRequest{ID: key, MType: "counter", Delta: &value})
	}
	jsonMetric, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("sendMetrics => %w", err)
	}
	_, err = zb.Write(jsonMetric)
	if err != nil {
		return fmt.Errorf("sendMetrics => %w", err)
	}
	if err = zb.Close(); err != nil {
		return fmt.Errorf("sendMetrics => %w", err)
	}
	r := client.NewRequest()
	r.Header.Set("Content-Encoding", "gzip")
	r.SetBody(buf)
	if _, err := r.Post(url); err != nil {
		return fmt.Errorf("sendMetrics => %w", err)
	}
	return nil
}
