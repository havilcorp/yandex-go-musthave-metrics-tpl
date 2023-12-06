package main

import (
	"net/http"
	"regexp"
	"strconv"
)

type MemStorage struct {
	listCounter map[string]float64
	listGauge   map[string]float64
}

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

	// marketName := tegexPath[1]
	marketVal := tegexPath[2]

	_, err := strconv.ParseFloat(marketVal, 64)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusOK)

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

	// marketName := tegexPath[1]
	marketVal := tegexPath[2]

	_, err := strconv.ParseFloat(marketVal, 64)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusOK)

}

func BadTypeHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusBadRequest)
}

func UpdateOtherHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusNotFound)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/update/counter/", UpdateCounterHandler)
	mux.HandleFunc("/update/gauge/", UpdateGaugeHandler)
	mux.HandleFunc("/update/", BadTypeHandler)
	mux.HandleFunc("/", UpdateOtherHandler)

	http.ListenAndServe(":8080", mux)
}
