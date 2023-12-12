package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/storage/memstorage"
)

var store = memstorage.MemStorage{Gauge: map[string]float64{}, Counter: map[string]int64{}}

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

func GetCounterMetricHandler(rw http.ResponseWriter, r *http.Request) {
	marketName := chi.URLParam(r, "name")
	if val, err := store.GetCounter(marketName); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	} else {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(fmt.Sprintf("%d", val)))
	}
}

func GetGaugeMetricHandler(rw http.ResponseWriter, r *http.Request) {
	marketName := chi.URLParam(r, "name")
	if val, err := store.GetGauge(marketName); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	} else {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(fmt.Sprintf("%g", val)))
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
