// Package main сервер для принятия и хранения метрик
package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/domain"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/config/server"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/handlers/rest"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/middleware"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/repositories"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/repositories/metric"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sirupsen/logrus"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

// main входная точка запуска сервера
func main() {
	conf := server.NewServerConfig()
	if err := conf.WriteByFlag(); err != nil {
		log.Fatal(err)
		return
	}
	if err := conf.WriteByEnv(); err != nil {
		log.Fatal(err)
		return
	}

	logrus.Infof("Build version: %s\n", buildVersion)
	logrus.Infof("Build date: %s\n", buildDate)
	logrus.Infof("Build commit: %s\n", buildCommit)

	provider := "memory"
	var db *sql.DB
	var err error

	if conf.DBConnect != "" {
		provider = "psql"
		db, err = sql.Open("pgx", conf.DBConnect)
		if err != nil {
			logrus.Errorf("pgx init => %v", err)
			return
		}
		defer func() {
			if err = db.Close(); err != nil {
				logrus.Error(err)
			}
		}()
		for _, sec := range []int{1, 3, 5} {
			err = db.PingContext(context.Background())
			if errors.Is(err, syscall.ECONNREFUSED) {
				time.Sleep(time.Duration(sec) * time.Second)
			} else {
				break
			}
		}
		if err != nil {
			logrus.Errorf("pgx init => %v", err)
			return
		}
	} else if conf.FileStoragePath != "" {
		provider = "file"
	}

	storageRepo := repositories.NewStorage(conf.FileStoragePath)

	metricFactory, err := metric.MetricFactory(provider, conf, db)
	if err != nil {
		logrus.Error(err)
		return
	}

	r := chi.NewRouter()
	if conf.TrustedSubnet != "" {
		r.Use(middleware.TrustedSubnetMiddleware(conf.TrustedSubnet))
	}
	r.Use(middleware.LogMiddleware)
	r.Use(middleware.GzipMiddleware)
	if conf.CryptoKey != "" {
		r.Use(middleware.RSAMiddleware(conf.CryptoKey))
	}
	r.Use(middleware.HashSHA256Middleware(conf.Key))

	rest.NewPPROFHandler().Register(r)
	rest.NewMainHandler(metricFactory).Register(r)
	rest.NewMetricHandler(metricFactory).Register(r)
	if conf.DBConnect != "" {
		dbRepository := repositories.NewDataBase(db)
		rest.NewPingHandler(dbRepository).Register(r)
	}

	var timeTicker *time.Ticker
	server := &http.Server{Addr: conf.ServerAddress, Handler: r}

	go func() {
		logrus.Infof("Сервер запушен %s", conf.ServerAddress)
		if err = server.ListenAndServe(); err != nil {
			logrus.Error(err)
		}
	}()

	if conf.StoreInterval != 0 {
		timeTicker = time.NewTicker(time.Second * time.Duration(conf.StoreInterval))
		go func() {
			for range timeTicker.C {
				var listGauges map[string]float64
				listGauges, err = metricFactory.GetAllGauge(context.Background())
				if err != nil {
					logrus.Error(err)
				}
				var listCounters map[string]int64
				listCounters, err = metricFactory.GetAllCounters(context.Background())
				if err != nil {
					logrus.Error(err)
				}
				metric := domain.Metric{
					Gauge:   listGauges,
					Counter: listCounters,
				}
				if err = storageRepo.SaveToFile(metric); err != nil {
					logrus.Error(err)
				}
			}
		}()
	}

	terminateSignals := make(chan os.Signal, 1)
	signal.Notify(terminateSignals, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-terminateSignals
	if err = server.Shutdown(context.Background()); err != nil {
		logrus.Error(err)
	}
	if timeTicker != nil {
		timeTicker.Stop()
	}
	listGauges, err := metricFactory.GetAllGauge(context.Background())
	if err != nil {
		logrus.Error(err)
	}
	listCounters, err := metricFactory.GetAllCounters(context.Background())
	if err != nil {
		logrus.Error(err)
	}
	metric := domain.Metric{
		Gauge:   listGauges,
		Counter: listCounters,
	}
	if err := storageRepo.SaveToFile(metric); err != nil {
		logrus.Error(err)
	}
	logrus.Info("Сервер остановлен")
}
