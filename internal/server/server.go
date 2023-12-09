package server

import (
	"net/http"

	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/handlers"
)

func CreateServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("/update/counter/", handlers.UpdateCounterHandler)
	mux.HandleFunc("/update/gauge/", handlers.UpdateGaugeHandler)
	mux.HandleFunc("/update/", handlers.BadTypeHandler)
	mux.HandleFunc("/print/", handlers.MainHandler)
	mux.HandleFunc("/", handlers.UpdateOtherHandler)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
