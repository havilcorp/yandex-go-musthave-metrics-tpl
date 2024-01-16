package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/config"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/handlers"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/middleware"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/storage/memstorage"
	"github.com/sirupsen/logrus"
)

func StartServer() error {

	var serverAddress string
	var storeInterval int
	var fileStoragePath string
	var isRestore bool

	config.WriteServerConfig(&serverAddress, &storeInterval, &fileStoragePath, &isRestore)

	logrus.Infof("StoreInterval: %d", storeInterval)
	logrus.Infof("FileStoragePath: %s", fileStoragePath)
	logrus.Infof("IsRestore: %t", isRestore)

	r := chi.NewRouter()
	server := &http.Server{Addr: serverAddress, Handler: r}

	store := *memstorage.NewMemStorage(storeInterval == 0)
	store.SetWfiteFileName(fileStoragePath)
	if isRestore {
		if err := store.LoadFromFile(); err != nil {
			logrus.Info(err)
			return err
		}
	}
	handlers.SetStore(store)

	// r.Use(middleware.Timeout(60 * time.Second))

	r.Use(middleware.LogMiddleware)
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

	go func() {
		logrus.Infof("Starting server on %s", serverAddress)
		if err := server.ListenAndServe(); err != nil {
			logrus.Info(err)
		}
	}()

	var timeTicker *time.Ticker

	if storeInterval != 0 {
		timeTicker = time.NewTicker(time.Second * time.Duration(storeInterval))
		go func() {
			for range timeTicker.C {
				if err := store.SaveToFile(); err != nil {
					logrus.Info(err)
				}
			}
		}()
	}

	terminateSignals := make(chan os.Signal, 1)
	signal.Notify(terminateSignals, syscall.SIGINT)
	<-terminateSignals
	server.Shutdown(context.Background())
	if timeTicker != nil {
		timeTicker.Stop()
	}
	if err := store.SaveToFile(); err != nil {
		logrus.Info(err)
		return err
	}
	logrus.Info("Сервер остановлен")
	return nil
}
