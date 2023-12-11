package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/handlers"
)

// var flagRunAddr string

// func parseFlags() {
// 	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
// 	flag.Parse()
// }

func CreateServer() {
	r := chi.NewRouter()

	// r.Use(middleware.RequestID)
	// r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", handlers.MainPageHandler)

	r.Route("/value", func(r chi.Router) {
		r.Get("/counter/{name}", handlers.GetCounterMetricHandler)
		r.Get("/gauge/{name}", handlers.GetGaugeMetricHandler)
	})

	r.Route("/update", func(r chi.Router) {
		r.Post("/counter/{name}/{value}", handlers.UpdateCounterHandler)
		r.Post("/gauge/{name}/{value}", handlers.UpdateGaugeHandler)
		r.Post("/{all}/{name}/{value}", handlers.BadRequestHandler)
	})

	// parseFlags()

	// fmt.Println("Running server on", flagRunAddr)
	// if err := http.ListenAndServe(flagRunAddr, r); err != nil {
	// 	panic(err)
	// }

	fmt.Println("Running server on", "8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}
