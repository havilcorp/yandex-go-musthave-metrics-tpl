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
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/middleware"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/storage"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/storage/file"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/storage/memory"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/storage/postgresql"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/transport/rest/metricupdate"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/transport/rest/metricvalue"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/transport/rest/other"
	"github.com/sirupsen/logrus"
)

func StartServer() error {

	conf := config.Config{}
	conf.WriteServerConfig()

	logrus.Infof("StoreInterval: %d", conf.StoreInterval)
	logrus.Infof("FileStoragePath: %s", conf.FileStoragePath)
	logrus.Infof("IsRestore: %t", conf.IsRestore)
	logrus.Infof("DBConnect: %s", conf.DBConnect)
	logrus.Infof("Key: %s", conf.Key)

	var storePtr storage.IStorage

	if conf.DBConnect != "" {
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

	if err := storePtr.Init(); err != nil {
		panic(err)
	}
	defer storePtr.Close()

	r := chi.NewRouter()

	// r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.LogMiddleware)
	r.Use(middleware.GzipMiddleware)
	r.Use(middleware.HashSHA256Middleware(conf.Key))

	metricupdate.NewHandler(storePtr).Register(r)
	metricvalue.NewHandler(storePtr).Register(r)
	other.NewHandler(storePtr).Register(r)

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
				if err := storePtr.SaveToFile(context.Background()); err != nil {
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
	if err := storePtr.SaveToFile(context.Background()); err != nil {
		logrus.Info(err)
		return err
	}
	logrus.Info("Приложение остановлено")
	return nil
}
