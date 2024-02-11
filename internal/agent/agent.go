package agent

import (
	"errors"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/config"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/metric"
	"github.com/sirupsen/logrus"
)

func workerSendeRequest(jobs <-chan metric.Metric, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		var err error
		for _, sec := range []int{1, 3, 5} {
			err = job.Send()
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

func StartAgent() {
	conf := config.NewConfig()
	if err := conf.WriteAgentConfig(); err != nil {
		logrus.Error(err)
		return
	}
	logrus.Info(conf)

	numJobs := runtime.GOMAXPROCS(0)
	numPool := conf.RateLimit

	jobs := make(chan metric.Metric, numJobs)

	var wg sync.WaitGroup
	for i := 0; i < numPool; i++ {
		wg.Add(1)
		go workerSendeRequest(jobs, &wg)
	}

	chDone := make(chan struct{})
	defer close(chDone)

	m := metric.NewMetric(conf)

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
	signal.Notify(terminateSignals, syscall.SIGINT)
	<-terminateSignals
	timePoolTracker.Stop()
	timeReportTracker.Stop()
	chDone <- struct{}{}
	close(jobs)
	wg.Wait()
	logrus.Info("Агент остановлен")
}
