package mertic

import (
	"fmt"
	"math/rand"
	"runtime"

	"github.com/go-resty/resty/v2"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/storage/memstorage"
)

var memStats runtime.MemStats

func WriteMetric(ms memstorage.Repositories) {
	runtime.ReadMemStats(&memStats)

	ms.AddGauge("Alloc", float64(memStats.Alloc))
	ms.AddGauge("BuckHashSys", float64(memStats.BuckHashSys))
	ms.AddGauge("Frees", float64(memStats.Frees))
	ms.AddGauge("GCCPUFraction", float64(memStats.GCCPUFraction))
	ms.AddGauge("GCSys", float64(memStats.GCSys))
	ms.AddGauge("HeapAlloc", float64(memStats.HeapAlloc))
	ms.AddGauge("HeapIdle", float64(memStats.HeapIdle))
	ms.AddGauge("HeapInuse", float64(memStats.HeapInuse))
	ms.AddGauge("HeapObjects", float64(memStats.HeapObjects))
	ms.AddGauge("HeapReleased", float64(memStats.HeapReleased))
	ms.AddGauge("HeapSys", float64(memStats.HeapSys))
	ms.AddGauge("LastGC", float64(memStats.LastGC))
	ms.AddGauge("Lookups", float64(memStats.Lookups))
	ms.AddGauge("MCacheInuse", float64(memStats.MCacheInuse))
	ms.AddGauge("MSpanSys", float64(memStats.MSpanSys))
	ms.AddGauge("Mallocs", float64(memStats.Mallocs))
	ms.AddGauge("NextGC", float64(memStats.NextGC))
	ms.AddGauge("NumForcedGC", float64(memStats.NumForcedGC))
	ms.AddGauge("NumGC", float64(memStats.NumGC))
	ms.AddGauge("OtherSys", float64(memStats.OtherSys))
	ms.AddGauge("PauseTotalNs", float64(memStats.PauseTotalNs))
	ms.AddGauge("StackInuse", float64(memStats.StackInuse))
	ms.AddGauge("StackSys", float64(memStats.StackSys))
	ms.AddGauge("Sys", float64(memStats.Sys))
	ms.AddGauge("TotalAlloc", float64(memStats.TotalAlloc))
	ms.AddCounter("PollCount", int64(1))
	ms.AddGauge("RandomValue", float64(rand.Intn(10)))
}

func SendMetric(address string, ms memstorage.Repositories) error {
	client := resty.New()

	gauges := ms.GetAllGauge()
	counters := ms.GetAllCounters()

	for key, val := range gauges {
		_, err := client.R().SetPathParams(map[string]string{
			"name":    key,
			"value":   fmt.Sprintf("%f", val),
			"address": address,
		}).Post("http://{address}/update/gauge/{name}/{value}")
		if err != nil {
			return err
		}
	}

	for key, val := range counters {
		_, err := client.R().SetPathParams(map[string]string{
			"name":    key,
			"value":   fmt.Sprintf("%d", val),
			"address": address,
		}).Post("http://{address}/update/counter/{name}/{value}")
		if err != nil {
			return err
		}
	}
	return nil
}
