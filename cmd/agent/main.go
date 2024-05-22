// Package main агент для отправки метрик на сервер
package main

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/config/agent"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/metric"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/middleware"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/havilcorp/yandex-go-musthave-metrics-tpl/pkg/proto/metric"
)

func workerSenderRequest(jobs <-chan metric.Metric, wg *sync.WaitGroup, conf *agent.Config) {
	defer wg.Done()
	for job := range jobs {
		var err error
		sender := job.Send
		if conf.AddressGRPC != "" {
			sender = job.SendByGRPC
		}
		for _, sec := range []int{1, 3, 5} {
			err = sender()
			if errors.Is(err, syscall.ECONNREFUSED) {
				time.Sleep(time.Duration(sec) * time.Second)
			} else {
				break
			}
		}
		if err != nil {
			logrus.Error(err)
		}
	}
}

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

// main входная точка запуска агента
func main() {
	conf := agent.NewAgentConfig()
	conf.WriteByFlag()
	if err := conf.WriteByEnv(); err != nil {
		log.Fatal(err)
	}

	logrus.Infof("Build version: %s\n", buildVersion)
	logrus.Infof("Build date: %s\n", buildDate)
	logrus.Infof("Build commit: %s\n", buildCommit)

	m := metric.NewMetric(conf)

	if conf.AddressGRPC != "" {
		cred := insecure.NewCredentials()
		if conf.CryptoCrt != "" {
			var err error
			cred, err = credentials.NewClientTLSFromFile(conf.CryptoCrt, "")
			if err != nil {
				logrus.Error(err)
				return
			}
		}
		conn, err := grpc.Dial(conf.AddressGRPC, grpc.WithTransportCredentials(cred), grpc.WithUnaryInterceptor(middleware.ClientInterceptor))
		if err != nil {
			logrus.Fatal(err)
		}
		defer func() {
			if err = conn.Close(); err != nil {
				logrus.Error(err)
			}
		}()
		client := pb.NewMetricClient(conn)
		m.AddMetricClient(client)
	}

	numJobs := runtime.GOMAXPROCS(0)
	numPool := conf.RateLimit

	jobs := make(chan metric.Metric, numJobs)

	var wg sync.WaitGroup
	for i := 0; i < numPool; i++ {
		wg.Add(1)
		go workerSenderRequest(jobs, &wg, conf)
	}

	chDone := make(chan struct{})
	defer close(chDone)

	timePoolTracker := time.NewTicker(time.Duration(conf.PollInterval) * time.Second)
	go func() {
		for {
			select {
			case <-timePoolTracker.C:
				go m.WriteMain()
				go m.WriteGopsutil()
			case <-chDone:
				return
			}
		}
	}()

	timeReportTracker := time.NewTicker(time.Duration(conf.ReportInterval) * time.Second)
	go func() {
		for {
			select {
			case <-timeReportTracker.C:
				jobs <- *m
			case <-chDone:
				return
			}
		}
	}()

	terminateSignals := make(chan os.Signal, 1)
	signal.Notify(terminateSignals, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-terminateSignals
	timePoolTracker.Stop()
	timeReportTracker.Stop()
	close(jobs)
	wg.Wait()
	logrus.Info("Агент остановлен")
}
