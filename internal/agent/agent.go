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

func StartAgent() error {
	config.WriteAgentConfig(&serverAddress, &reportInterval, &pollInterval)

	// timerInterval := time.NewTicker(time.Duration(reportInterval) * time.Second)
	// timerPool := time.NewTicker(time.Duration(pollInterval) * time.Second)
	// timerPool := time.NewTicker(time.Duration(1) * time.Second)

	// for v := range ch {
	// 	// do some stuff
	// }

	stopTimer := make(chan bool)
	timeTicker := time.NewTicker(time.Second)
	go func() {
		i := 0
		defer func() { stopTimer <- true }()
		for {
			select {
			case <-timeTicker.C:
				if i%pollInterval == 0 {
					mertic.WriteMetric(store)
				}
				if i%reportInterval == 0 {
					mertic.SendMetric(serverAddress, store)
				}
				i++
			case <-stopTimer:
				return
			}
		}
	}()

	terminateSignals := make(chan os.Signal, 1)
	signal.Notify(terminateSignals, syscall.SIGINT)
	<-terminateSignals
	timeTicker.Stop()
	stopTimer <- true
	logrus.Info("Агент остановлен")
	return nil
}
