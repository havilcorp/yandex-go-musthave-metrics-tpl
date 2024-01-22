package agent

import (
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
				mertic.WriteMetric(store)
			}
			if i%conf.ReportInterval == 0 {
				mertic.SendMetric(conf.ServerAddress, store)
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
