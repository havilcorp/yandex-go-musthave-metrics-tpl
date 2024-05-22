// Package metric вспомогательный пакет для работы агента
package metric

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/hmac"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	mr "math/rand"
	"net"
	"runtime"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/domain"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/config/agent"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/cryptorsa"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/sirupsen/logrus"

	pb "github.com/havilcorp/yandex-go-musthave-metrics-tpl/pkg/proto/metric"
)

type Metric struct {
	mutex        *sync.Mutex
	config       *agent.Config
	value        map[string]float64
	delta        map[string]int64
	metricClient pb.MetricClient
}

func NewMetric(config *agent.Config) *Metric {
	return &Metric{
		mutex:  &sync.Mutex{},
		config: config,
		value:  make(map[string]float64, 0),
		delta:  make(map[string]int64, 0),
	}
}

func (m *Metric) AddMetricClient(mc pb.MetricClient) {
	m.metricClient = mc
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

// WriteGopsutil получение дополнительных метрик
//
//   - TotalMemory
//   - FreeMemory
//   - CPUutilization
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

// WriteMain получение основных метрик
//
//   - Alloc
//   - BuckHashSys
//   - Frees
//   - GCCPUFraction
//   - GCSys
//   - HeapAlloc
//   - HeapIdle
//   - HeapInuse
//   - HeapObjects
//   - HeapReleased
//   - HeapSys
//   - LastGC
//   - Lookups
//   - MCacheInuse
//   - MCacheSys
//   - MSpanInuse
//   - MSpanSys
//   - Mallocs
//   - NextGC
//   - NumForcedGC
//   - NumGC
//   - OtherSys
//   - PauseTotalNs
//   - StackInuse
//   - StackSys
//   - Sys
//   - TotalAlloc
//   - RandomValue
//   - PollCount
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
	m.value["RandomValue"] = float64(mr.Intn(10))
	m.delta["PollCount"] = int64(1)
}

func (m *Metric) SendByGRPC() error {
	gauge := make([]*pb.Gauge, 0)
	counter := make([]*pb.Counter, 0)
	for key, value := range m.value {
		gauge = append(gauge, &pb.Gauge{
			Key:   key,
			Value: value,
		})
	}
	for key, delta := range m.delta {
		counter = append(counter, &pb.Counter{
			Key:   key,
			Value: delta,
		})
	}
	if m.metricClient == nil {
		return errors.New("metricClient has a nil pointer")
	}
	_, err := m.metricClient.UpdateMetricBulk(context.Background(), &pb.UpdateMetricBulkRequest{
		Gauge:   gauge,
		Counter: counter,
	})
	return err
}

// Send отправка метрик на сервер
func (m *Metric) Send() error {
	logrus.Info("SEND")
	client := resty.New()
	url := fmt.Sprintf("http://%s/updates", m.config.ServerAddress)
	buf := bytes.NewBuffer(nil)
	zb := gzip.NewWriter(buf)
	metrics := make([]domain.MetricRequest, 0)
	for key, value := range m.value {
		metrics = append(metrics, domain.MetricRequest{ID: key, MType: "gauge", Value: &value})
	}
	for key, delta := range m.delta {
		metrics = append(metrics, domain.MetricRequest{ID: key, MType: "counter", Delta: &delta})
	}
	jsonMetric, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("sendMetrics => %w", err)
	}
	if m.config.CryptoKey != "" {
		var pub *rsa.PublicKey
		pub, err = cryptorsa.LoadPublicKey(m.config.CryptoKey)
		if err != nil {
			return fmt.Errorf("LoadPublicKey: %w", err)
		}
		jsonMetric, err = cryptorsa.EncryptOAEP(pub, jsonMetric)
		if err != nil {
			return fmt.Errorf("EncryptOAEP: %w", err)
		}
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
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return fmt.Errorf("get ip address => %w", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			logrus.Info(err)
		}
	}()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	r.Header.Set("X-Real-IP", localAddr.IP.To4().String())
	if _, err := r.Post(url); err != nil {
		return fmt.Errorf("sendMetrics => %w", err)
	}
	return nil
}
