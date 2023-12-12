package agent

import (
	"fmt"
	"time"

	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/config"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/mertic"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/storage/memstorage"
)

var store = memstorage.MemStorage{Gauge: map[string]float64{}, Counter: map[string]int64{}}

var serverAddress string
var reportInterval int
var pollInterval int

func StartAgent() error {
	config.WriteAgentConfig(&serverAddress, &reportInterval, &pollInterval)

	timer := time.NewTicker(1 * time.Second)
	stop := make(chan bool)

	go func() {
		defer func() { stop <- true }()
		i := 0
		for {
			select {
			case <-timer.C:
				time.Sleep(1 * time.Second)

				if i%pollInterval == 0 {
					mertic.WriteMetric(store)
				}

				if i%reportInterval == 0 && i != 0 {
					mertic.SendMetric(serverAddress, store)
					i = 0
				}

				i++

			case <-stop:
				fmt.Println("Закрытие горутины")
				return
			}
		}
	}()

	return nil
}
