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
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/storage"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/storage/file"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/storage/memory"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/storage/postgresql"
	"github.com/sirupsen/logrus"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func StartServer() error {

	conf := config.Config{}
	conf.WriteServerConfig()

	logrus.Infof("StoreInterval: %d", conf.StoreInterval)
	logrus.Infof("FileStoragePath: %s", conf.FileStoragePath)
	logrus.Infof("IsRestore: %t", conf.IsRestore)
	logrus.Infof("DbConnect: %s", conf.DbConnect)

	var storePtr storage.IStorage

	if conf.DbConnect != "" {
		logrus.Info("PsqlStorage")
		storePtr = &postgresql.PsqlStorage{
			Conf: &conf,
		}
	} else if conf.FileStoragePath != "" {
		logrus.Info("FileStorage")
		storePtr = &file.FileStorage{
			Conf:    &conf,
			Gauge:   map[string]float64{},
			Counter: map[string]int64{},
		}
	} else {
		logrus.Info("MemStorage")
		storePtr = &memory.MemStorage{
			Gauge:   map[string]float64{},
			Counter: map[string]int64{},
		}
	}

	if err := storePtr.Init(context.Background()); err != nil {
		logrus.Info(err)
		return nil
	}
	defer storePtr.Close()

	handlers.SetStore(storePtr)

	r := chi.NewRouter()

	// r.Use(middleware.Timeout(60 * time.Second))

	r.Use(middleware.LogMiddleware)
	r.Use(middleware.GzipMiddleware)

	r.Get("/", handlers.MainPageHandler)

	r.Get("/ping", handlers.CheckDBHandler)

	r.Route("/value", func(r chi.Router) {
		r.Post("/", handlers.GetMetricHandler)
		r.Get("/counter/{name}", handlers.GetCounterMetricHandler)
		r.Get("/gauge/{name}", handlers.GetGaugeMetricHandler)
	})

	r.Route("/updates", func(r chi.Router) {
		r.Post("/", handlers.UpdateBulkHandler)
	})

	r.Route("/update", func(r chi.Router) {
		r.Post("/", handlers.UpdateHandler)
		r.Post("/counter/{name}/{value}", handlers.UpdateCounterHandler)
		r.Post("/gauge/{name}/{value}", handlers.UpdateGaugeHandler)
		r.Post("/{all}/{name}/{value}", handlers.BadRequestHandler)
	})

	var timeTicker *time.Ticker
	server := &http.Server{Addr: conf.ServerAddress, Handler: r}

	go func() {
		logrus.Infof("Starting server on %s", conf.ServerAddress)
		if err := server.ListenAndServe(); err != nil {
			logrus.Info(err)
		}
	}()

	if conf.StoreInterval != 0 {
		timeTicker = time.NewTicker(time.Second * time.Duration(conf.StoreInterval))
		go func() {
			for range timeTicker.C {
				if err := storePtr.SaveToFile(); err != nil {
					logrus.Info(err)
				}
			}
		}()
	}

	terminateSignals := make(chan os.Signal, 1)
	signal.Notify(terminateSignals, syscall.SIGINT)
	<-terminateSignals
	if err := server.Shutdown(context.Background()); err != nil {
		logrus.Info(err)
	}
	if timeTicker != nil {
		timeTicker.Stop()
	}
	if err := storePtr.SaveToFile(); err != nil {
		logrus.Info(err)
		return err
	}
	logrus.Info("Приложение остановлено")
	return nil
}
