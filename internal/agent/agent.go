package agent

import (
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/config"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/mertic"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/storage/memory"
	"github.com/sirupsen/logrus"
)

func StartAgent() {
	conf := config.Config{}
	conf.WriteAgentConfig()

	store := memory.MemStorage{Gauge: map[string]float64{}, Counter: map[string]int64{}}

	timeTicker := time.NewTicker(time.Second)
	go func() {
		i := 0
		for range timeTicker.C {
			if i%conf.PollInterval == 0 {
				if err := mertic.WriteMetric(store); err != nil {
					logrus.Info(err)
					panic(err)
				}
			}
			if i%conf.ReportInterval == 0 {
				var err error
				for _, sec := range []int{1, 3, 5} {
					err = mertic.SendMetric(conf.ServerAddress, store)
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
			i++
		}
	}()

	terminateSignals := make(chan os.Signal, 1)
	signal.Notify(terminateSignals, syscall.SIGINT)
	<-terminateSignals
	timeTicker.Stop()
	logrus.Info("Агент остановлен")
}
