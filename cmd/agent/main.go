package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/storage"
)

var store = storage.MemStorage{Gauge: map[string]float64{}, Counter: map[string]int64{}}

func main() {

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	i := 0
	for {
		time.Sleep(1 * time.Second)

		if i%2 == 0 {
			store.AddGauge("Alloc", float64(memStats.Alloc))
			store.AddGauge("BuckHashSys", float64(memStats.BuckHashSys))
			store.AddGauge("Frees", float64(memStats.Frees))
			store.AddGauge("GCCPUFraction", float64(memStats.GCCPUFraction))
			store.AddGauge("GCSys", float64(memStats.GCSys))
			store.AddGauge("HeapAlloc", float64(memStats.HeapAlloc))
			store.AddGauge("HeapIdle", float64(memStats.HeapIdle))
			store.AddGauge("HeapInuse", float64(memStats.HeapInuse))
			store.AddGauge("HeapObjects", float64(memStats.HeapObjects))
			store.AddGauge("HeapReleased", float64(memStats.HeapReleased))
			store.AddGauge("HeapSys", float64(memStats.HeapSys))
			store.AddGauge("LastGC", float64(memStats.LastGC))
			store.AddGauge("Lookups", float64(memStats.Lookups))
			store.AddGauge("MCacheInuse", float64(memStats.MCacheInuse))
			store.AddGauge("MSpanSys", float64(memStats.MSpanSys))
			store.AddGauge("Mallocs", float64(memStats.Mallocs))
			store.AddGauge("NextGC", float64(memStats.NextGC))
			store.AddGauge("NumForcedGC", float64(memStats.NumForcedGC))
			store.AddGauge("NumGC", float64(memStats.NumGC))
			store.AddGauge("OtherSys", float64(memStats.OtherSys))
			store.AddGauge("PauseTotalNs", float64(memStats.PauseTotalNs))
			store.AddGauge("StackInuse", float64(memStats.StackInuse))
			store.AddGauge("StackSys", float64(memStats.StackSys))
			store.AddGauge("Sys", float64(memStats.Sys))
			store.AddGauge("TotalAlloc", float64(memStats.TotalAlloc))
			store.AddCounter("PollCount", int64(1))
			store.AddGauge("RandomValue", float64(rand.Intn(10)))
		}

		if i%10 == 0 && i != 0 {
			for key, val := range store.Gauge {
				url := fmt.Sprintf("http://127.0.0.1:8080/update/gauge/%s/%f", key, val)
				_, err := http.Post(url, "text/plain", nil)
				if err != nil {
					panic(err)
				}
			}

			for key, val := range store.Counter {
				url := fmt.Sprintf("http://127.0.0.1:8080/update/counter/%s/%d", key, val)
				_, err := http.Post(url, "text/plain", nil)
				if err != nil {
					panic(err)
				}
			}

		}

		i++
	}
}
