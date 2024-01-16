package mertic

import (
	"bytes"
	"compress/gzip"
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
	ms.AddGauge("MCacheSys", float64(memStats.MCacheSys))
	ms.AddGauge("MSpanInuse", float64(memStats.MSpanInuse))
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
	ms.AddGauge("RandomValue", float64(rand.Intn(10)))
	ms.AddCounter("PollCount", int64(1))
}

func SendMetric(address string, ms memstorage.Repositories) error {
	client := resty.New()

	gauges := ms.GetAllGauge()
	counters := ms.GetAllCounters()

	for key, val := range gauges {
		url := fmt.Sprintf("http://%s/update", address)
		buf := bytes.NewBuffer(nil)
		zb := gzip.NewWriter(buf)
		_, err := zb.Write([]byte(fmt.Sprintf(`{"id":"%s","type":"gauge","value":%f}`, key, val)))
		if err != nil {
			return err
		}
		if err = zb.Close(); err != nil {
			return err
		}
		r := client.NewRequest()
		r.Header.Set("Content-Encoding", "gzip")
		r.SetBody(buf)
		if _, err = r.Post(url); err != nil {
			return err
		}
	}

	for key, val := range counters {
		url := fmt.Sprintf("http://%s/update", address)
		buf := bytes.NewBuffer(nil)
		zb := gzip.NewWriter(buf)
		_, err := zb.Write([]byte(fmt.Sprintf(`{"id":"%s","type":"counter","delta":%d}`, key, val)))
		if err != nil {
			return err
		}
		if err = zb.Close(); err != nil {
			return err
		}
		r := client.NewRequest()
		r.Header.Set("Content-Encoding", "gzip")
		r.SetBody(buf)
		if _, err = r.Post(url); err != nil {
			return err
		}
	}
	return nil
}
