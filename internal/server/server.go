package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/config"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/handlers"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/logger"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/middleware"
	"go.uber.org/zap"
)

func StartServer() error {

	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}
	defer zapLogger.Sync()
	sugar := *zapLogger.Sugar()

	r := chi.NewRouter()

	// r.Use(middleware.Timeout(60 * time.Second))

	r.Use(logger.WithLogging)
	r.Use(middleware.GzipMiddleware)

	r.Get("/", handlers.MainPageHandler)

	r.Route("/value", func(r chi.Router) {
		r.Post("/", handlers.GetMetricHandler)
		r.Get("/counter/{name}", handlers.GetCounterMetricHandler)
		r.Get("/gauge/{name}", handlers.GetGaugeMetricHandler)
	})

	r.Route("/update", func(r chi.Router) {
		r.Post("/", handlers.UpdateHandler)
		r.Post("/counter/{name}/{value}", handlers.UpdateCounterHandler)
		r.Post("/gauge/{name}/{value}", handlers.UpdateGaugeHandler)
		r.Post("/{all}/{name}/{value}", handlers.BadRequestHandler)
	})

	var serverAddress string

	config.WriteServerConfig(&serverAddress)

	sugar.Infow(
		"Starting server",
		"addr", serverAddress,
	)

	return http.ListenAndServe(serverAddress, r)
}
