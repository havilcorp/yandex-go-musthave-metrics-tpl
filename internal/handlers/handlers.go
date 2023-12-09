package handlers

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/storage"
)

var store = storage.MemStorage{Gauge: map[string]float64{}, Counter: map[string]int64{}}

func UpdateCounterHandler(rw http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	regex := `^\/update\/counter\/([a-zA-Z0-9]+)\/([a-zA-Z0-9.]+)$`
	tegexPath := regexp.MustCompile(regex).FindStringSubmatch(r.URL.Path)

	if len(tegexPath) != 3 {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	marketName := tegexPath[1]
	marketVal := tegexPath[2]

	val64, err := strconv.ParseInt(marketVal, 0, 64)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusOK)

	if err := store.AddCounter(marketName, val64); err != nil {
		panic(err)
	}

}

func UpdateGaugeHandler(rw http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	regex := `^\/update\/gauge\/([a-zA-Z0-9]+)\/([a-zA-Z0-9.]+)$`
	tegexPath := regexp.MustCompile(regex).FindStringSubmatch(r.URL.Path)

	if len(tegexPath) != 3 {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	marketName := tegexPath[1]
	marketVal := tegexPath[2]

	val64, err := strconv.ParseFloat(marketVal, 64)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := store.AddGauge(marketName, val64); err != nil {
		panic(err)
	}

	rw.WriteHeader(http.StatusOK)

}

func BadTypeHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusBadRequest)
}

func UpdateOtherHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusNotFound)
}

func MainHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte(fmt.Sprint(store.Gauge) + fmt.Sprint(store.Counter)))
}
