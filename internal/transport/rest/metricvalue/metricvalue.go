package metricvalue

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/models"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/storage"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/transport/rest"
	"github.com/sirupsen/logrus"
)

type handler struct {
	store storage.IStorage
}

func NewHandler(store storage.IStorage) rest.Handler {
	return &handler{store: store}
}

func (h *handler) Register(router *chi.Mux) {
	router.Route("/value", func(r chi.Router) {
		r.Post("/", h.GetMetricHandler)
		r.Get("/counter/{name}", h.GetCounterMetricHandler)
		r.Get("/gauge/{name}", h.GetGaugeMetricHandler)
	})
}

func (h *handler) GetMetricHandler(rw http.ResponseWriter, r *http.Request) {
	var req models.MetricsRequest
	rw.Header().Set("Content-Type", "application/json")
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		logrus.Error(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	if req.MType == models.TypeMetricsCounter {
		val, err := h.store.GetCounter(r.Context(), req.ID)
		if err != nil {
			if errors.Is(err, storage.ErrValueNotFound) {
				rw.WriteHeader(http.StatusNotFound)
				return
			}
			logrus.Error(err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		rw.WriteHeader(http.StatusOK)
		resp := models.MetricsRequest{
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
	}
	if req.MType == models.TypeMetricsGauge {
		val, err := h.store.GetGauge(r.Context(), req.ID)
		if err != nil {
			if errors.Is(err, storage.ErrValueNotFound) {
				rw.WriteHeader(http.StatusNotFound)
				return
			}
			logrus.Error(err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		rw.WriteHeader(http.StatusOK)
		resp := models.MetricsRequest{
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
	}
}
func (h *handler) GetCounterMetricHandler(rw http.ResponseWriter, r *http.Request) {
	marketName := chi.URLParam(r, "name")
	val, err := h.store.GetCounter(r.Context(), marketName)
	if err != nil {
		if errors.Is(err, storage.ErrValueNotFound) {
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
func (h *handler) GetGaugeMetricHandler(rw http.ResponseWriter, r *http.Request) {
	marketName := chi.URLParam(r, "name")
	val, err := h.store.GetGauge(r.Context(), marketName)
	if err != nil {
		if errors.Is(err, storage.ErrValueNotFound) {
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
