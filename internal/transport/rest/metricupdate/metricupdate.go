package metricupdate

import (
	"encoding/json"
	"net/http"
	"strconv"

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
	router.Route("/updates", func(r chi.Router) {
		r.Post("/", h.UpdateBulkHandler)
	})
	router.Route("/update", func(r chi.Router) {
		r.Post("/", h.UpdateHandler)
		r.Post("/counter/{name}/{value}", h.UpdateCounterHandler)
		r.Post("/gauge/{name}/{value}", h.UpdateGaugeHandler)
		r.Post("/{all}/{name}/{value}", h.BadRequestHandler)
	})
}

func (h *handler) UpdateBulkHandler(rw http.ResponseWriter, r *http.Request) {
	metrics := make([]models.MetricsRequest, 0)
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&metrics); err != nil {
		logrus.Error(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	gauge := make([]models.GaugeModel, 0)
	counter := make([]models.CounterModel, 0)
	for _, m := range metrics {
		if m.MType == models.TypeMetricsGauge {
			gauge = append(gauge, models.GaugeModel{Key: m.ID, Value: *m.Value})
		} else if m.MType == models.TypeMetricsCounter {
			counter = append(counter, models.CounterModel{Key: m.ID, Value: *m.Delta})
		}
	}
	if err := h.store.AddGaugeBulk(r.Context(), gauge); err != nil {
		logrus.Error(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := h.store.AddCounterBulk(r.Context(), counter); err != nil {
		logrus.Error(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *handler) UpdateHandler(rw http.ResponseWriter, r *http.Request) {
	var req models.MetricsRequest
	rw.Header().Set("Content-Type", "application/json")
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		logrus.Error(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	if req.MType == models.TypeMetricsCounter {
		if err := h.store.AddCounter(r.Context(), req.ID, *req.Delta); err != nil {
			logrus.Error(err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		val, err := h.store.GetCounter(r.Context(), req.ID)
		if err != nil {
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
		if err := h.store.AddGauge(r.Context(), req.ID, *req.Value); err != nil {
			logrus.Error(err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		val, err := h.store.GetGauge(r.Context(), req.ID)
		if err != nil {
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

func (h *handler) UpdateCounterHandler(rw http.ResponseWriter, r *http.Request) {
	marketName := chi.URLParam(r, "name")
	marketVal := chi.URLParam(r, "value")
	marketValInt64, err := strconv.ParseInt(marketVal, 0, 64)
	if err != nil {
		logrus.Error(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.store.AddCounter(r.Context(), marketName, marketValInt64); err != nil {
		logrus.Error(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusOK)
}
func (h *handler) UpdateGaugeHandler(rw http.ResponseWriter, r *http.Request) {
	marketName := chi.URLParam(r, "name")
	marketVal := chi.URLParam(r, "value")
	marketValFloat64, err := strconv.ParseFloat(marketVal, 64)
	if err != nil {
		logrus.Error(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.store.AddGauge(r.Context(), marketName, marketValFloat64); err != nil {
		logrus.Error(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusOK)
}
func (h *handler) BadRequestHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusBadRequest)
}
