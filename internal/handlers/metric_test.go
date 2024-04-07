package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/domain"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/config"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/repositories/metric"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/mocks"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMetricHandler_UpdateBulkHandler(t *testing.T) {
	metricHandler := mocks.NewIMetric(t)

	metricHandler.On("AddGaugeBulk", mock.Anything, []domain.Gauge{{Key: "GAUGE", Value: float64(1.1)}}).Return(nil)
	metricHandler.On("AddCounterBulk", mock.Anything, []domain.Counter{{Key: "COUNTER", Value: int64(1)}}).Return(nil)

	delta := int64(1)
	value := float64(1.1)

	type args struct {
		statusCode int
		data       []domain.MetricRequest
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "UpdateBulkHandler",
			args: args{
				statusCode: 200,
				data: []domain.MetricRequest{{
					ID:    "COUNTER",
					MType: "counter",
					Delta: &delta,
				}, {
					ID:    "GAUGE",
					MType: "gauge",
					Value: &value,
				}},
			},
		},
		{
			name: "UpdateBulkHandler",
			args: args{
				statusCode: 400,
				data: []domain.MetricRequest{{
					ID:    "COUNTER",
					MType: "none",
					Delta: &delta,
				}, {
					ID:    "GAUGE",
					MType: "none",
					Value: &value,
				}},
			},
		},
	}
	for _, tt := range tests {
		jsonData, err := json.Marshal(tt.args.data)
		if err != nil {
			t.Errorf("jsonData %v", err)
		}
		r := httptest.NewRequest(http.MethodPost, "/updates", strings.NewReader(string(jsonData)))
		rw := httptest.NewRecorder()
		t.Run(tt.name, func(t *testing.T) {
			h := NewMetricHandler(metricHandler)
			h.UpdateBulkHandler(rw, r)
			res := rw.Result()
			assert.Equal(t, tt.args.statusCode, res.StatusCode)
			defer res.Body.Close()
		})
	}
}

func TestMetricHandler_UpdateHandler(t *testing.T) {
	metricHandler := mocks.NewIMetric(t)

	metricHandler.On("AddGauge", mock.Anything, "GAUGE", float64(1.1)).Return(nil)
	metricHandler.On("AddCounter", mock.Anything, "COUNTER", int64(1)).Return(nil)
	metricHandler.On("GetGauge", mock.Anything, "GAUGE").Return(float64(1.1), nil)
	metricHandler.On("GetCounter", mock.Anything, "COUNTER").Return(int64(1), nil)

	delta := int64(1)
	value := float64(1.1)

	type args struct {
		statusCode int
		data       domain.MetricRequest
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "UpdateHandler 1",
			args: args{
				statusCode: 200,
				data: domain.MetricRequest{
					ID:    "COUNTER",
					MType: "counter",
					Delta: &delta,
				},
			},
		},
		{
			name: "UpdateHandler 2",
			args: args{
				statusCode: 200,
				data: domain.MetricRequest{
					ID:    "GAUGE",
					MType: "gauge",
					Value: &value,
				},
			},
		},
		{
			name: "UpdateHandler 3",
			args: args{
				statusCode: 400,
				data: domain.MetricRequest{
					ID:    "GAUGE",
					MType: "none",
					Value: &value,
				},
			},
		},
	}
	for _, tt := range tests {
		jsonData, err := json.Marshal(tt.args.data)
		if err != nil {
			t.Errorf("jsonData %v", err)
		}
		r := httptest.NewRequest(http.MethodPost, "/updates", strings.NewReader(string(jsonData)))
		rw := httptest.NewRecorder()
		t.Run(tt.name, func(t *testing.T) {
			h := NewMetricHandler(metricHandler)
			h.UpdateHandler(rw, r)
			res := rw.Result()
			assert.Equal(t, tt.args.statusCode, res.StatusCode)
			defer res.Body.Close()
			// data, err := io.ReadAll(res.Body)
			// if err != nil {
			// 	t.Errorf("expected error to be nil got %v", err)
			// }
			// if string(data) == "" {
			// 	t.Errorf("expected ABC got %v", string(data))
			// }
		})
	}
}

func TestMetricHandler_UpdateCounterHandler(t *testing.T) {
	metricHandler := mocks.NewIMetric(t)

	metricHandler.On("AddCounter", mock.Anything, "COUNTER", int64(1)).Return(nil)

	type args struct {
		statusCode int
		key        string
		value      string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "UpdateCounterHandler",
			args: args{
				statusCode: 200,
				key:        "COUNTER",
				value:      "1",
			},
		},
		{
			name: "UpdateCounterHandler2",
			args: args{
				statusCode: 400,
				key:        "COUNTER",
				value:      "1.2",
			},
		},
	}
	for _, tt := range tests {
		r := httptest.NewRequest(http.MethodPost, "/update/counter/{name}/{value}", nil)
		rw := httptest.NewRecorder()
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("name", tt.args.key)
		rctx.URLParams.Add("value", tt.args.value)
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		t.Run(tt.name, func(t *testing.T) {
			h := NewMetricHandler(metricHandler)
			h.UpdateCounterHandler(rw, r)
			res := rw.Result()
			assert.Equal(t, tt.args.statusCode, res.StatusCode)
			defer res.Body.Close()
			// data, err := io.ReadAll(res.Body)
			// if err != nil {
			// 	t.Errorf("expected error to be nil got %v", err)
			// }
			// if string(data) == "" {
			// 	t.Errorf("expected ABC got %v", string(data))
			// }
		})
	}
}

func TestMetricHandler_UpdateGaugeHandler(t *testing.T) {
	metricHandler := mocks.NewIMetric(t)

	metricHandler.On("AddGauge", mock.Anything, "GAUGE", float64(1.1)).Return(nil)

	type args struct {
		statusCode int
		key        string
		value      string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "UpdateGaugeHandler",
			args: args{
				statusCode: 200,
				key:        "GAUGE",
				value:      "1.1",
			},
		},
		{
			name: "UpdateGaugeHandler 2",
			args: args{
				statusCode: 400,
				key:        "GAUGE",
				value:      "asd",
			},
		},
	}
	for _, tt := range tests {
		r := httptest.NewRequest(http.MethodPost, "/update/gauge/{name}/{value}", nil)
		rw := httptest.NewRecorder()
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("name", tt.args.key)
		rctx.URLParams.Add("value", tt.args.value)
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		t.Run(tt.name, func(t *testing.T) {
			h := NewMetricHandler(metricHandler)
			h.UpdateGaugeHandler(rw, r)
			res := rw.Result()
			assert.Equal(t, tt.args.statusCode, res.StatusCode)
			defer res.Body.Close()
			// data, err := io.ReadAll(res.Body)
			// if err != nil {
			// 	t.Errorf("expected error to be nil got %v", err)
			// }
			// if string(data) == "" {
			// 	t.Errorf("expected ABC got %v", string(data))
			// }
		})
	}
}

func TestMetricHandler_BadRequestHandler(t *testing.T) {
	metricHandler := mocks.NewIMetric(t)
	type args struct {
		statusCode int
		all        string
		key        string
		value      string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "BadRequestHandler",
			args: args{
				statusCode: 400,
				all:        "all",
				key:        "key",
				value:      "1",
			},
		},
	}
	for _, tt := range tests {
		r := httptest.NewRequest(http.MethodPost, "/update/{all}/{name}/{value}", nil)
		rw := httptest.NewRecorder()
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("all", tt.args.all)
		rctx.URLParams.Add("name", tt.args.key)
		rctx.URLParams.Add("value", tt.args.value)
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		t.Run(tt.name, func(t *testing.T) {
			h := NewMetricHandler(metricHandler)
			h.BadRequestHandler(rw, r)
			res := rw.Result()
			assert.Equal(t, tt.args.statusCode, res.StatusCode)
			defer res.Body.Close()
			// data, err := io.ReadAll(res.Body)
			// if err != nil {
			// 	t.Errorf("expected error to be nil got %v", err)
			// }
			// if string(data) == "" {
			// 	t.Errorf("expected ABC got %v", string(data))
			// }
		})
	}
}

func TestMetricHandler_GetMetricHandler(t *testing.T) {
	metricHandler := mocks.NewIMetric(t)
	metricHandler.On("GetCounter", mock.Anything, "COUNTER").Return(int64(1), nil)
	metricHandler.On("GetGauge", mock.Anything, "GAUGE").Return(float64(1.1), nil)
	type args struct {
		statusCode int
		data       domain.MetricRequest
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "GetMetricHandler 1",
			args: args{
				statusCode: 200,
				data: domain.MetricRequest{
					ID:    "COUNTER",
					MType: "counter",
				},
			},
		},
		{
			name: "GetMetricHandler 2",
			args: args{
				statusCode: 200,
				data: domain.MetricRequest{
					ID:    "GAUGE",
					MType: "gauge",
				},
			},
		},
		{
			name: "GetMetricHandler 3",
			args: args{
				statusCode: 400,
				data: domain.MetricRequest{
					ID:    "GAUGE",
					MType: "none",
				},
			},
		},
	}
	for _, tt := range tests {
		jsonData, err := json.Marshal(tt.args.data)
		if err != nil {
			t.Errorf("jsonData %v", err)
		}
		r := httptest.NewRequest(http.MethodPost, "/value", strings.NewReader(string(jsonData)))
		rw := httptest.NewRecorder()
		t.Run(tt.name, func(t *testing.T) {
			h := NewMetricHandler(metricHandler)
			h.GetMetricHandler(rw, r)
			res := rw.Result()
			assert.Equal(t, tt.args.statusCode, res.StatusCode)
			defer res.Body.Close()
			// data, err := io.ReadAll(res.Body)
			// if res.StatusCode == 200 {
			// 	if err != nil {
			// 		t.Errorf("expected error to be nil got %v", err)
			// 	}
			// 	if string(data) == "" {
			// 		t.Errorf("expected ABC got %v", string(data))
			// 	}
			// }
		})
	}
}

func TestMetricHandler_GetCounterMetricHandler(t *testing.T) {
	metricHandler := mocks.NewIMetric(t)
	metricHandler.On("GetCounter", mock.Anything, "COUNTER").Return(int64(1), nil)
	type args struct {
		statusCode int
		name       string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "GetCounterMetricHandler",
			args: args{
				statusCode: 200,
				name:       "COUNTER",
			},
		},
	}
	for _, tt := range tests {
		r := httptest.NewRequest(http.MethodPost, "/value/counter/{name}", nil)
		rw := httptest.NewRecorder()
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("name", tt.args.name)
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		t.Run(tt.name, func(t *testing.T) {
			h := NewMetricHandler(metricHandler)
			h.GetCounterMetricHandler(rw, r)
			res := rw.Result()
			assert.Equal(t, tt.args.statusCode, res.StatusCode)
			defer res.Body.Close()
			// data, err := io.ReadAll(res.Body)
			// if err != nil {
			// 	t.Errorf("expected error to be nil got %v", err)
			// }
			// if string(data) == "" {
			// 	t.Errorf("expected ABC got %v", string(data))
			// }
		})
	}
}

func TestMetricHandler_GetGaugeMetricHandler(t *testing.T) {
	metricHandler := mocks.NewIMetric(t)
	metricHandler.On("GetGauge", mock.Anything, "GAUGE").Return(float64(1.1), nil)
	type args struct {
		statusCode int
		name       string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "GetGaugeMetricHandler",
			args: args{
				statusCode: 200,
				name:       "GAUGE",
			},
		},
	}
	for _, tt := range tests {
		r := httptest.NewRequest(http.MethodPost, "/value/gauge/{name}", nil)
		rw := httptest.NewRecorder()
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("name", tt.args.name)
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		t.Run(tt.name, func(t *testing.T) {
			h := NewMetricHandler(metricHandler)
			h.GetGaugeMetricHandler(rw, r)
			res := rw.Result()
			assert.Equal(t, tt.args.statusCode, res.StatusCode)
			defer res.Body.Close()
			// data, err := io.ReadAll(res.Body)
			// if err != nil {
			// 	t.Errorf("expected error to be nil got %v", err)
			// }
			// if string(data) == "" {
			// 	t.Errorf("expected ABC got %v", string(data))
			// }
		})
	}
}

func Example() {
	conf := config.NewConfig()
	err := conf.WriteServerConfig()
	if err != nil {
		logrus.Error(err)
		return
	}
	metricFactory, err := metric.MetricFactory("memory", conf, nil)
	if err != nil {
		logrus.Error(err)
		return
	}
	r := chi.NewRouter()
	NewMetricHandler(metricFactory).Register(r)
}
