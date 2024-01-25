package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/models"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/storage"
	"github.com/sirupsen/logrus"
)

var store storage.IStorage

func SetStore(s storage.IStorage) {
	store = s
}

func CheckDBHandler(rw http.ResponseWriter, r *http.Request) {
	if err := store.Ping(); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	} else {
		rw.WriteHeader(http.StatusOK)
	}
}

func UpdateHandler(rw http.ResponseWriter, r *http.Request) {
	var req models.MetricsRequest
	rw.Header().Set("Content-Type", "application/json")
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		logrus.Info(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	if req.MType == models.TypeMetricsCounter {
		if err := store.AddCounter(r.Context(), req.ID, *req.Delta); err != nil {
			logrus.Info(err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		if val, ok := store.GetCounter(r.Context(), req.ID); ok {
			rw.WriteHeader(http.StatusOK)
			resp := models.MetricsRequest{
				ID:    req.ID,
				MType: req.MType,
				Delta: &val,
			}
			enc := json.NewEncoder(rw)
			if err := enc.Encode(resp); err != nil {
				logrus.Info(err)
				return
			}
		} else {
			rw.WriteHeader(http.StatusNotFound)
		}
	}
	if req.MType == models.TypeMetricsGauge {
		if err := store.AddGauge(r.Context(), req.ID, *req.Value); err != nil {
			logrus.Info(err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		if val, ok := store.GetGauge(r.Context(), req.ID); ok {
			rw.WriteHeader(http.StatusOK)
			resp := models.MetricsRequest{
				ID:    req.ID,
				MType: req.MType,
				Value: &val,
			}
			enc := json.NewEncoder(rw)
			if err := enc.Encode(resp); err != nil {
				logrus.Info(err)
				return
			}
		} else {
			rw.WriteHeader(http.StatusNotFound)
		}
	}

}

func UpdateBulkHandler(rw http.ResponseWriter, r *http.Request) {
	metrics := make([]models.MetricsRequest, 0)
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&metrics); err != nil {
		logrus.Info(err)
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
	if err := store.AddGaugeBulk(r.Context(), gauge); err != nil {
		logrus.Info(err)
	}
	if err := store.AddCounterBulk(r.Context(), counter); err != nil {
		logrus.Info(err)
	}
}

func UpdateCounterHandler(rw http.ResponseWriter, r *http.Request) {

	marketName := chi.URLParam(r, "name")
	marketVal := chi.URLParam(r, "value")

	marketValInt64, err := strconv.ParseInt(marketVal, 0, 64)
	if err != nil {
		logrus.Info(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := store.AddCounter(r.Context(), marketName, marketValInt64); err != nil {
		logrus.Info(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func UpdateGaugeHandler(rw http.ResponseWriter, r *http.Request) {

	marketName := chi.URLParam(r, "name")
	marketVal := chi.URLParam(r, "value")

	marketValFloat64, err := strconv.ParseFloat(marketVal, 64)
	if err != nil {
		logrus.Info(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := store.AddGauge(r.Context(), marketName, marketValFloat64); err != nil {
		logrus.Info(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusOK)

}

func GetMetricHandler(rw http.ResponseWriter, r *http.Request) {
	var req models.MetricsRequest
	rw.Header().Set("Content-Type", "application/json")
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		logrus.Info(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	if req.MType == models.TypeMetricsCounter {
		if val, ok := store.GetCounter(r.Context(), req.ID); ok {
			rw.WriteHeader(http.StatusOK)
			resp := models.MetricsRequest{
				ID:    req.ID,
				MType: req.MType,
				Delta: &val,
			}
			enc := json.NewEncoder(rw)
			if err := enc.Encode(resp); err != nil {
				logrus.Info(err)
				return
			}
		} else {
			rw.WriteHeader(http.StatusNotFound)
		}
	}
	if req.MType == models.TypeMetricsGauge {
		if val, ok := store.GetGauge(r.Context(), req.ID); ok {
			rw.WriteHeader(http.StatusOK)
			resp := models.MetricsRequest{
				ID:    req.ID,
				MType: req.MType,
				Value: &val,
			}
			enc := json.NewEncoder(rw)
			if err := enc.Encode(resp); err != nil {
				logrus.Info(err)
				return
			}
		} else {
			rw.WriteHeader(http.StatusNotFound)
		}
	}
}

func GetCounterMetricHandler(rw http.ResponseWriter, r *http.Request) {
	marketName := chi.URLParam(r, "name")
	if val, ok := store.GetCounter(r.Context(), marketName); ok {
		rw.WriteHeader(http.StatusOK)
		_, err := rw.Write([]byte(fmt.Sprintf("%d", val)))
		if err != nil {
			logrus.Info(err)
		}
	} else {
		rw.WriteHeader(http.StatusNotFound)
	}
}

func GetGaugeMetricHandler(rw http.ResponseWriter, r *http.Request) {
	marketName := chi.URLParam(r, "name")
	if val, ok := store.GetGauge(r.Context(), marketName); ok {
		rw.WriteHeader(http.StatusOK)
		_, err := rw.Write([]byte(fmt.Sprintf("%g", val)))
		if err != nil {
			logrus.Info(err)
		}
	} else {
		rw.WriteHeader(http.StatusNotFound)
	}
}

func MainPageHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "text/html")
	liCounter := ""
	for key, item := range store.GetAllCounters(r.Context()) {
		liCounter += fmt.Sprintf("<li>%s: %d</li>", key, item)
	}
	liGauge := ""
	for key, item := range store.GetAllGauge(r.Context()) {
		liGauge += fmt.Sprintf("<li>%s: %f</li>", key, item)
	}
	html := fmt.Sprintf(`
	<html>
		<body>
			<br/>
			Counters
			<br/>
			<ul>%s</ul>
			<br/>
			Gauges
			<br/>
			<ul>%s</ul>
		</body>
	</html>
	`, liCounter, liGauge)
	_, err := rw.Write([]byte(html))
	if err != nil {
		logrus.Info(err)
	}
}

func BadRequestHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusBadRequest)
}
