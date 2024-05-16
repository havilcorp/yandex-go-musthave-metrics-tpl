// Package rest роуты сервера
package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/domain"
	"github.com/sirupsen/logrus"
)

type IMetric interface {
	AddGauge(ctx context.Context, key string, gauge float64) error
	AddCounter(ctx context.Context, key string, counter int64) error
	AddGaugeBulk(ctx context.Context, list []domain.Gauge) error
	AddCounterBulk(ctx context.Context, list []domain.Counter) error
	GetCounter(ctx context.Context, key string) (int64, error)
	GetGauge(ctx context.Context, key string) (float64, error)
}

type MetricHandler struct {
	metric IMetric
}

// NewMetricHandler инициализация хендлера
func NewMetricHandler(metric IMetric) *MetricHandler {
	return &MetricHandler{
		metric: metric,
	}
}

// Register регистрация роутов
func (h *MetricHandler) Register(router *chi.Mux) {
	router.Route("/updates", func(r chi.Router) {
		r.Post("/", h.UpdateBulkHandler)
	})
	router.Route("/update", func(r chi.Router) {
		r.Post("/", h.UpdateHandler)
		r.Post("/counter/{name}/{value}", h.UpdateCounterHandler)
		r.Post("/gauge/{name}/{value}", h.UpdateGaugeHandler)
		r.Post("/{all}/{name}/{value}", h.BadRequestHandler)
	})
	router.Route("/value", func(r chi.Router) {
		r.Post("/", h.GetMetricHandler)
		r.Get("/counter/{name}", h.GetCounterMetricHandler)
		r.Get("/gauge/{name}", h.GetGaugeMetricHandler)
	})
}

// UpdateBulkHandler хендлер обновления метрик пачкой
func (h *MetricHandler) UpdateBulkHandler(rw http.ResponseWriter, r *http.Request) {
	metrics := make([]domain.MetricRequest, 0)
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&metrics); err != nil {
		logrus.Error(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	gauge := make([]domain.Gauge, 0)
	counter := make([]domain.Counter, 0)
	for _, m := range metrics {
		if m.MType == domain.TypeMetricsGauge {
			gauge = append(gauge, domain.Gauge{Key: m.ID, Value: *m.Value})
		} else if m.MType == domain.TypeMetricsCounter {
			counter = append(counter, domain.Counter{Key: m.ID, Value: *m.Delta})
		} else {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	if err := h.metric.AddGaugeBulk(r.Context(), gauge); err != nil {
		logrus.Error(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := h.metric.AddCounterBulk(r.Context(), counter); err != nil {
		logrus.Error(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// UpdateHandler хендлер обновления одной из метрик
func (h *MetricHandler) UpdateHandler(rw http.ResponseWriter, r *http.Request) {
	var req domain.MetricRequest
	rw.Header().Set("Content-Type", "application/json")
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		logrus.Error(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	if req.MType == domain.TypeMetricsCounter {
		if err := h.metric.AddCounter(r.Context(), req.ID, *req.Delta); err != nil {
			logrus.Error(err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		val, err := h.metric.GetCounter(r.Context(), req.ID)
		if err != nil {
			logrus.Error(err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		rw.WriteHeader(http.StatusOK)
		resp := domain.MetricRequest{
			ID:    req.ID,
			MType: req.MType,
			Delta: &val,
		}
		enc := json.NewEncoder(rw)
		if err := enc.Encode(resp); err != nil {
			logrus.Error(err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else if req.MType == domain.TypeMetricsGauge {
		if err := h.metric.AddGauge(r.Context(), req.ID, *req.Value); err != nil {
			logrus.Error(err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		val, err := h.metric.GetGauge(r.Context(), req.ID)
		if err != nil {
			logrus.Error(err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		rw.WriteHeader(http.StatusOK)
		resp := domain.MetricRequest{
			ID:    req.ID,
			MType: req.MType,
			Value: &val,
		}
		enc := json.NewEncoder(rw)
		if err := enc.Encode(resp); err != nil {
			logrus.Error(err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
}

// UpdateCounterHandler хендлер обновления метрики
func (h *MetricHandler) UpdateCounterHandler(rw http.ResponseWriter, r *http.Request) {
	marketName := chi.URLParam(r, "name")
	marketVal := chi.URLParam(r, "value")
	marketValInt64, err := strconv.ParseInt(marketVal, 0, 64)
	if err != nil {
		logrus.Error(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.metric.AddCounter(r.Context(), marketName, marketValInt64); err != nil {
		logrus.Error(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusOK)
}

// UpdateGaugeHandler хендлер обновления метрики
func (h *MetricHandler) UpdateGaugeHandler(rw http.ResponseWriter, r *http.Request) {
	marketName := chi.URLParam(r, "name")
	marketVal := chi.URLParam(r, "value")
	marketValFloat64, err := strconv.ParseFloat(marketVal, 64)
	if err != nil {
		logrus.Error(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.metric.AddGauge(r.Context(), marketName, marketValFloat64); err != nil {
		logrus.Error(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusOK)
}

func (h *MetricHandler) BadRequestHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusBadRequest)
}

// GetMetricHandler хендлер выдачи значения метрик
func (h *MetricHandler) GetMetricHandler(rw http.ResponseWriter, r *http.Request) {
	var req domain.MetricRequest
	rw.Header().Set("Content-Type", "application/json")
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		logrus.Error(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	if req.MType == domain.TypeMetricsCounter {
		val, err := h.metric.GetCounter(r.Context(), req.ID)
		if err != nil {
			if errors.Is(err, domain.ErrValueNotFound) {
				rw.WriteHeader(http.StatusNotFound)
				return
			}
			logrus.Error(err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		rw.WriteHeader(http.StatusOK)
		resp := domain.MetricRequest{
			ID:    req.ID,
			MType: req.MType,
			Delta: &val,
		}
		enc := json.NewEncoder(rw)
		if err := enc.Encode(resp); err != nil {
			logrus.Error(err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else if req.MType == domain.TypeMetricsGauge {
		val, err := h.metric.GetGauge(r.Context(), req.ID)
		if err != nil {
			if errors.Is(err, domain.ErrValueNotFound) {
				rw.WriteHeader(http.StatusNotFound)
				return
			}
			logrus.Error(err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		rw.WriteHeader(http.StatusOK)
		resp := domain.MetricRequest{
			ID:    req.ID,
			MType: req.MType,
			Value: &val,
		}
		enc := json.NewEncoder(rw)
		if err := enc.Encode(resp); err != nil {
			logrus.Error(err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
}

// GetCounterMetricHandler хендлер выдачи значения метрики
func (h *MetricHandler) GetCounterMetricHandler(rw http.ResponseWriter, r *http.Request) {
	marketName := chi.URLParam(r, "name")
	val, err := h.metric.GetCounter(r.Context(), marketName)
	if err != nil {
		if errors.Is(err, domain.ErrValueNotFound) {
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		logrus.Error(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
	_, err = rw.Write([]byte(fmt.Sprintf("%d", val)))
	if err != nil {
		logrus.Error(err)
		return
	}
}

// GetGaugeMetricHandler хендлер выдачи значения метрики
func (h *MetricHandler) GetGaugeMetricHandler(rw http.ResponseWriter, r *http.Request) {
	marketName := chi.URLParam(r, "name")
	val, err := h.metric.GetGauge(r.Context(), marketName)
	if err != nil {
		if errors.Is(err, domain.ErrValueNotFound) {
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		logrus.Error(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
	_, err = rw.Write([]byte(fmt.Sprintf("%g", val)))
	if err != nil {
		logrus.Error(err)
		return
	}
}
