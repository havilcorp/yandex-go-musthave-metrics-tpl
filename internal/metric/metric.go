package metric

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/config"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/models"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

type Metric struct {
	mutex  *sync.Mutex
	config *config.Config
	value  map[string]float64
	delta  map[string]int64
}

func NewMetric(config *config.Config) *Metric {
	return &Metric{
		mutex:  &sync.Mutex{},
		config: config,
		value:  make(map[string]float64, 0),
		delta:  make(map[string]int64, 0),
	}
}

func (m *Metric) String() string {
	out := ""
	for key, value := range m.value {
		out += fmt.Sprintf("[%s]: %f\n", key, value)
	}
	for key, delta := range m.delta {
		out += fmt.Sprintf("[%s]: %d\n", key, delta)
	}
	return out
}

func (m *Metric) WriteGopsutil() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	v, _ := mem.VirtualMemory()
	m.value["TotalMemory"] = float64(v.Total)
	m.value["FreeMemory"] = float64(v.Free)
	c, err := cpu.Percent(4*time.Second, true)
	if err != nil {
		panic(err)
	}
	for i, val := range c {
		m.value[fmt.Sprintf("CPUutilization%d", i+1)] = float64(val)
	}
}

func (m *Metric) WriteMain() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	m.value["Alloc"] = float64(memStats.Alloc)
	m.value["BuckHashSys"] = float64(memStats.BuckHashSys)
	m.value["Frees"] = float64(memStats.Frees)
	m.value["GCCPUFraction"] = float64(memStats.GCCPUFraction)
	m.value["GCSys"] = float64(memStats.GCSys)
	m.value["HeapAlloc"] = float64(memStats.HeapAlloc)
	m.value["HeapIdle"] = float64(memStats.HeapIdle)
	m.value["HeapInuse"] = float64(memStats.HeapInuse)
	m.value["HeapObjects"] = float64(memStats.HeapObjects)
	m.value["HeapReleased"] = float64(memStats.HeapReleased)
	m.value["HeapSys"] = float64(memStats.HeapSys)
	m.value["LastGC"] = float64(memStats.LastGC)
	m.value["Lookups"] = float64(memStats.Lookups)
	m.value["MCacheInuse"] = float64(memStats.MCacheInuse)
	m.value["MCacheSys"] = float64(memStats.MCacheSys)
	m.value["MSpanInuse"] = float64(memStats.MSpanInuse)
	m.value["MSpanSys"] = float64(memStats.MSpanSys)
	m.value["Mallocs"] = float64(memStats.Mallocs)
	m.value["NextGC"] = float64(memStats.NextGC)
	m.value["NumForcedGC"] = float64(memStats.NumForcedGC)
	m.value["NumGC"] = float64(memStats.NumGC)
	m.value["OtherSys"] = float64(memStats.OtherSys)
	m.value["PauseTotalNs"] = float64(memStats.PauseTotalNs)
	m.value["StackInuse"] = float64(memStats.StackInuse)
	m.value["StackSys"] = float64(memStats.StackSys)
	m.value["Sys"] = float64(memStats.Sys)
	m.value["TotalAlloc"] = float64(memStats.TotalAlloc)
	m.value["RandomValue"] = float64(rand.Intn(10))
	m.delta["PollCount"] = int64(1)
}

func (m *Metric) Send() error {
	client := resty.New()
	url := fmt.Sprintf("http://%s/updates", m.config.ServerAddress)
	buf := bytes.NewBuffer(nil)
	zb := gzip.NewWriter(buf)
	metrics := make([]models.MetricsRequest, 0)
	for key, value := range m.value {
		metrics = append(metrics, models.MetricsRequest{ID: key, MType: "gauge", Value: &value})
	}
	for key, delta := range m.delta {
		metrics = append(metrics, models.MetricsRequest{ID: key, MType: "counter", Delta: &delta})
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
	if m.config.Key != "" {
		h := hmac.New(sha256.New, []byte(m.config.Key))
		h.Write(jsonMetric)
		hashSha256 := hex.EncodeToString(h.Sum(nil))
		r.Header.Set("HashSHA256", hashSha256)
	}
	r.SetBody(buf)
	if _, err := r.Post(url); err != nil {
		return fmt.Errorf("sendMetrics => %w", err)
	}
	return nil
}
