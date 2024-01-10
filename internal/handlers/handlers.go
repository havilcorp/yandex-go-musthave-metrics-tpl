package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/models"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/storage/memstorage"
)

var store = memstorage.MemStorage{Gauge: map[string]float64{}, Counter: map[string]int64{}}

func UpdateHandler(rw http.ResponseWriter, r *http.Request) {
	var req models.MetricsRequest
	rw.Header().Set("Content-Type", "application/json")
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	if req.MType == models.TypeMetricsCounter {
		fmt.Println("UpdateHandler", req.ID, req.MType, *req.Delta, 0)
		if err := store.AddCounter(req.ID, *req.Delta); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		if val, ok := store.GetCounter(req.ID); ok {
			rw.WriteHeader(http.StatusOK)
			resp := models.MetricsRequest{
				ID:    req.ID,
				MType: req.MType,
				Delta: &val,
			}
			enc := json.NewEncoder(rw)
			if err := enc.Encode(resp); err != nil {
				return
			}
		} else {
			rw.WriteHeader(http.StatusNotFound)
		}
	}
	if req.MType == models.TypeMetricsGauge {
		fmt.Println("UpdateHandler", req.ID, req.MType, 0, *req.Value)
		if err := store.AddGauge(req.ID, *req.Value); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		if val, ok := store.GetGauge(req.ID); ok {
			rw.WriteHeader(http.StatusOK)
			resp := models.MetricsRequest{
				ID:    req.ID,
				MType: req.MType,
				Value: &val,
			}
			enc := json.NewEncoder(rw)
			if err := enc.Encode(resp); err != nil {
				return
			}
		} else {
			rw.WriteHeader(http.StatusNotFound)
		}
	}

}

func UpdateCounterHandler(rw http.ResponseWriter, r *http.Request) {

	marketName := chi.URLParam(r, "name")
	marketVal := chi.URLParam(r, "value")

	marketValInt64, err := strconv.ParseInt(marketVal, 0, 64)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := store.AddCounter(marketName, marketValInt64); err != nil {
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
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := store.AddGauge(marketName, marketValFloat64); err != nil {
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
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	if req.MType == models.TypeMetricsCounter {
		fmt.Println("GetMetricHandler", req.ID, req.MType)
		if val, ok := store.GetCounter(req.ID); ok {
			rw.WriteHeader(http.StatusOK)
			resp := models.MetricsRequest{
				ID:    req.ID,
				MType: req.MType,
				Delta: &val,
			}
			enc := json.NewEncoder(rw)
			if err := enc.Encode(resp); err != nil {
				return
			}
		} else {
			rw.WriteHeader(http.StatusNotFound)
		}
	}
	if req.MType == models.TypeMetricsGauge {
		fmt.Println("GetMetricHandler", req.ID, req.MType)
		if val, ok := store.GetGauge(req.ID); ok {
			rw.WriteHeader(http.StatusOK)
			resp := models.MetricsRequest{
				ID:    req.ID,
				MType: req.MType,
				Value: &val,
			}
			enc := json.NewEncoder(rw)
			if err := enc.Encode(resp); err != nil {
				return
			}
		} else {
			rw.WriteHeader(http.StatusNotFound)
		}
	}
}

func GetCounterMetricHandler(rw http.ResponseWriter, r *http.Request) {
	marketName := chi.URLParam(r, "name")
	if val, ok := store.GetCounter(marketName); ok {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(fmt.Sprintf("%d", val)))
	} else {
		rw.WriteHeader(http.StatusNotFound)
	}
}

func GetGaugeMetricHandler(rw http.ResponseWriter, r *http.Request) {
	marketName := chi.URLParam(r, "name")
	if val, ok := store.GetGauge(marketName); ok {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(fmt.Sprintf("%g", val)))
	} else {
		rw.WriteHeader(http.StatusNotFound)
	}
}

func MainPageHandler(rw http.ResponseWriter, r *http.Request) {
	liCounter := ""
	for key, item := range store.Counter {
		liCounter += fmt.Sprintf("<li>%s: %d</li>", key, item)
	}
	liGauge := ""
	for key, item := range store.Gauge {
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
	</html>`,
		liCounter, liGauge)
	rw.Write([]byte(html))
}

func BadRequestHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusBadRequest)
}
