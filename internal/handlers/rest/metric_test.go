package rest

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/domain"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/config/server"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/repositories/metric"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/mocks"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMetricHandler_UpdateBulkHandler(t *testing.T) {
	metricHandler := mocks.NewMetricRouter(t)

	metricHandler.On("AddGaugeBulk", mock.Anything, []domain.Gauge{{Key: "GAUGE", Value: float64(1.1)}}).Return(nil)
	metricHandler.On("AddCounterBulk", mock.Anything, []domain.Counter{{Key: "COUNTER", Value: int64(1)}}).Return(nil)
	metricHandler.On("AddGaugeBulk", mock.Anything, []domain.Gauge{{Key: "GAUGE", Value: float64(2.1)}}).Return(errors.New(""))
	metricHandler.On("AddCounterBulk", mock.Anything, []domain.Counter{{Key: "COUNTER", Value: int64(2)}}).Return(errors.New(""))

	delta := int64(1)
	value := float64(1.1)
	deltaErr := int64(2)
	valueErr := float64(2.1)

	type args struct {
		data       []domain.MetricRequest
		statusCode int
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
		{
			name: "UpdateBulkHandler2",
			args: args{
				statusCode: 500,
				data: []domain.MetricRequest{{
					ID:    "COUNTER",
					MType: "counter",
					Delta: &deltaErr,
				}, {
					ID:    "GAUGE",
					MType: "gauge",
					Value: &value,
				}},
			},
		},
		{
			name: "UpdateBulkHandler3",
			args: args{
				statusCode: 500,
				data: []domain.MetricRequest{{
					ID:    "COUNTER",
					MType: "counter",
					Delta: &delta,
				}, {
					ID:    "GAUGE",
					MType: "gauge",
					Value: &valueErr,
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
			defer func() {
				if err := res.Body.Close(); err != nil {
					logrus.Info(err)
				}
			}()
		})
	}
}

func TestMetricHandler_UpdateHandler(t *testing.T) {
	metricHandler := mocks.NewMetricRouter(t)

	metricHandler.On("AddGauge", mock.Anything, "GAUGE", float64(1.1)).Return(nil)
	metricHandler.On("AddCounter", mock.Anything, "COUNTER", int64(1)).Return(nil)
	metricHandler.On("GetGauge", mock.Anything, "GAUGE").Return(float64(1.1), nil)
	metricHandler.On("GetCounter", mock.Anything, "COUNTER").Return(int64(1), nil)

	metricHandler.On("AddGauge", mock.Anything, "GAUGE", float64(2.1)).Return(errors.New(""))
	metricHandler.On("AddCounter", mock.Anything, "COUNTER", int64(2)).Return(errors.New(""))
	metricHandler.On("GetGauge", mock.Anything, "GAUGE").Return(float64(2.1), nil)
	metricHandler.On("GetCounter", mock.Anything, "COUNTER").Return(int64(2), nil)

	metricHandler.On("AddGauge", mock.Anything, "GAUGEErr", float64(3.1)).Return(nil)
	metricHandler.On("AddCounter", mock.Anything, "COUNTERErr", int64(3)).Return(nil)
	metricHandler.On("GetGauge", mock.Anything, "GAUGEErr").Return(float64(3.1), errors.New(""))
	metricHandler.On("GetCounter", mock.Anything, "COUNTERErr").Return(int64(3), errors.New(""))

	delta := int64(1)
	value := float64(1.1)
	deltaErrAdd := int64(2)
	valueErrAdd := float64(2.1)
	deltaErrGet := int64(3)
	valueErrGet := float64(3.1)

	type args struct {
		data       domain.MetricRequest
		statusCode int
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
		{
			name: "UpdateHandler 4",
			args: args{
				statusCode: 500,
				data: domain.MetricRequest{
					ID:    "GAUGE",
					MType: "gauge",
					Value: &valueErrAdd,
				},
			},
		},
		{
			name: "UpdateHandler 5",
			args: args{
				statusCode: 500,
				data: domain.MetricRequest{
					ID:    "COUNTER",
					MType: "counter",
					Delta: &deltaErrAdd,
				},
			},
		},
		{
			name: "UpdateHandler 6",
			args: args{
				statusCode: 500,
				data: domain.MetricRequest{
					ID:    "GAUGEErr",
					MType: "gauge",
					Value: &valueErrGet,
				},
			},
		},
		{
			name: "UpdateHandler 7",
			args: args{
				statusCode: 500,
				data: domain.MetricRequest{
					ID:    "COUNTERErr",
					MType: "counter",
					Delta: &deltaErrGet,
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
			defer func() {
				if err := res.Body.Close(); err != nil {
					logrus.Info(err)
				}
			}()
		})
	}
}

func TestMetricHandler_UpdateCounterHandler(t *testing.T) {
	metricHandler := mocks.NewMetricRouter(t)

	metricHandler.On("AddCounter", mock.Anything, "COUNTER", int64(1)).Return(nil)
	metricHandler.On("AddCounter", mock.Anything, "COUNTER", int64(2)).Return(errors.New(""))

	type args struct {
		key        string
		value      string
		statusCode int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "UpdateCounterHandler 1",
			args: args{
				statusCode: 200,
				key:        "COUNTER",
				value:      "1",
			},
		},
		{
			name: "UpdateCounterHandler 2",
			args: args{
				statusCode: 400,
				key:        "COUNTER",
				value:      "1.2",
			},
		},
		{
			name: "UpdateCounterHandler 3",
			args: args{
				statusCode: 400,
				key:        "COUNTER",
				value:      "2",
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
			defer func() {
				if err := res.Body.Close(); err != nil {
					logrus.Info(err)
				}
			}()
		})
	}
}

func TestMetricHandler_UpdateGaugeHandler(t *testing.T) {
	metricHandler := mocks.NewMetricRouter(t)

	metricHandler.On("AddGauge", mock.Anything, "GAUGE", float64(1.1)).Return(nil)
	metricHandler.On("AddGauge", mock.Anything, "GAUGE", float64(2.1)).Return(errors.New(""))

	type args struct {
		key        string
		value      string
		statusCode int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "UpdateGaugeHandler 1",
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
		{
			name: "UpdateGaugeHandler 3",
			args: args{
				statusCode: 400,
				key:        "GAUGE",
				value:      "2.1",
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
			defer func() {
				if err := res.Body.Close(); err != nil {
					logrus.Info(err)
				}
			}()
		})
	}
}

func TestMetricHandler_BadRequestHandler(t *testing.T) {
	metricHandler := mocks.NewMetricRouter(t)
	type args struct {
		all        string
		key        string
		value      string
		statusCode int
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
			defer func() {
				if err := res.Body.Close(); err != nil {
					logrus.Info(err)
				}
			}()
		})
	}
}

func TestMetricHandler_GetMetricHandler(t *testing.T) {
	metricHandler := mocks.NewMetricRouter(t)
	metricHandler.On("GetCounter", mock.Anything, "COUNTER").Return(int64(1), nil)
	metricHandler.On("GetGauge", mock.Anything, "GAUGE").Return(float64(1.1), nil)

	metricHandler.On("GetCounter", mock.Anything, "COUNTER404").Return(int64(2), domain.ErrValueNotFound)
	metricHandler.On("GetGauge", mock.Anything, "GAUGE404").Return(float64(2.1), domain.ErrValueNotFound)

	metricHandler.On("GetCounter", mock.Anything, "COUNTER500").Return(int64(3), errors.New(""))
	metricHandler.On("GetGauge", mock.Anything, "GAUGE500").Return(float64(3.1), errors.New(""))
	type args struct {
		data       domain.MetricRequest
		statusCode int
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
		{
			name: "GetMetricHandler 4",
			args: args{
				statusCode: 404,
				data: domain.MetricRequest{
					ID:    "COUNTER404",
					MType: "counter",
				},
			},
		},
		{
			name: "GetMetricHandler 5",
			args: args{
				statusCode: 404,
				data: domain.MetricRequest{
					ID:    "GAUGE404",
					MType: "gauge",
				},
			},
		},
		{
			name: "GetMetricHandler 6",
			args: args{
				statusCode: 500,
				data: domain.MetricRequest{
					ID:    "COUNTER500",
					MType: "counter",
				},
			},
		},
		{
			name: "GetMetricHandler 7",
			args: args{
				statusCode: 500,
				data: domain.MetricRequest{
					ID:    "GAUGE500",
					MType: "gauge",
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
			defer func() {
				if err := res.Body.Close(); err != nil {
					logrus.Info(err)
				}
			}()
		})
	}
}

func TestMetricHandler_GetCounterMetricHandler(t *testing.T) {
	metricHandler := mocks.NewMetricRouter(t)
	metricHandler.On("GetCounter", mock.Anything, "COUNTER").Return(int64(1), nil)
	metricHandler.On("GetCounter", mock.Anything, "COUNTER404").Return(int64(2), domain.ErrValueNotFound)
	metricHandler.On("GetCounter", mock.Anything, "COUNTER500").Return(int64(3), errors.New(""))
	type args struct {
		name       string
		statusCode int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "GetCounterMetricHandler 200",
			args: args{
				statusCode: 200,
				name:       "COUNTER",
			},
		},
		{
			name: "GetCounterMetricHandler 404",
			args: args{
				statusCode: 404,
				name:       "COUNTER404",
			},
		},
		{
			name: "GetCounterMetricHandler 500",
			args: args{
				statusCode: 500,
				name:       "COUNTER500",
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
			defer func() {
				if err := res.Body.Close(); err != nil {
					logrus.Info(err)
				}
			}()
		})
	}
}

func TestMetricHandler_GetGaugeMetricHandler(t *testing.T) {
	metricHandler := mocks.NewMetricRouter(t)
	metricHandler.On("GetGauge", mock.Anything, "GAUGE").Return(float64(1.1), nil)
	metricHandler.On("GetGauge", mock.Anything, "GAUGE404").Return(float64(2.1), domain.ErrValueNotFound)
	metricHandler.On("GetGauge", mock.Anything, "GAUGE500").Return(float64(3.1), errors.New(""))
	type args struct {
		name       string
		statusCode int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "GetGaugeMetricHandler 200",
			args: args{
				statusCode: 200,
				name:       "GAUGE",
			},
		},
		{
			name: "GetGaugeMetricHandler 404",
			args: args{
				statusCode: 404,
				name:       "GAUGE404",
			},
		},
		{
			name: "GetGaugeMetricHandler 500",
			args: args{
				statusCode: 500,
				name:       "GAUGE500",
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
			defer func() {
				if err := res.Body.Close(); err != nil {
					logrus.Info(err)
				}
			}()
		})
	}
}

func Example() {
	conf := server.NewServerConfig()
	conf.WriteByFlag()
	if err := conf.WriteByEnv(); err != nil {
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

func TestMetricHandler_Register(t *testing.T) {
	r := chi.NewRouter()
	metricHandler := mocks.NewMetricRouter(t)
	h := NewMetricHandler(metricHandler)
	h.Register(r)
}
