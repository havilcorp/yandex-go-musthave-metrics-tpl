package agent

import (
	"errors"
	"fmt"
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
			logrus.Info(err)
		}
	}
}

func StartAgent() {
	conf := config.NewConfig()
	if err := conf.WriteAgentConfig(); err != nil {
		logrus.Info(err)
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

	var mutex sync.Mutex
	m := metric.NewMetric(&mutex, conf)

	timeTracker := time.NewTicker(time.Second)
	defer timeTracker.Stop()
	go func(chDone chan struct{}) {
		i := 0
		for {
			select {
			case <-timeTracker.C:
				i++
				if i%conf.PollInterval == 0 {
					go m.WriteMain()
					go m.WriteGopsutil()
				}
				if i%conf.ReportInterval == 0 {
					i = 0
					fmt.Println(m)
					jobs <- *m
				}
			case <-chDone:
				return
			}
		}
	}(chDone)

	terminateSignals := make(chan os.Signal, 1)
	signal.Notify(terminateSignals, syscall.SIGINT)
	<-terminateSignals
	close(jobs)
	chDone <- struct{}{}
	wg.Wait()
	timeTracker.Stop()
	logrus.Info("Агент остановлен")
}
