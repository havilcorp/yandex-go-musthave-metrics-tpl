package agent

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/config"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/mertic"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/storage/memstorage"
	"github.com/sirupsen/logrus"
)

var store = memstorage.NewMemStorage(false)

var serverAddress string
var reportInterval int
var pollInterval int

func StartAgent() {
	config.WriteAgentConfig(&serverAddress, &reportInterval, &pollInterval)

	timeTicker := time.NewTicker(time.Second)
	go func() {
		i := 0
		for range timeTicker.C {
			if i%pollInterval == 0 {
				mertic.WriteMetric(store)
			}
			if i%reportInterval == 0 {
				mertic.SendMetric(serverAddress, store)
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
