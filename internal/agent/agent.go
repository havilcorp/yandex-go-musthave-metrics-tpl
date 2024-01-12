package agent

import (
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

	// timerInterval := time.NewTicker(time.Duration(reportInterval) * time.Second)
	// timerPool := time.NewTicker(time.Duration(pollInterval) * time.Second)
	timerPool := time.NewTicker(time.Duration(1) * time.Second)
	i := 0

	for {
		select {
		case <-timerPool.C:
			if i%pollInterval == 0 {
				mertic.WriteMetric(store)
			}
			if i%reportInterval == 0 {
				mertic.SendMetric(serverAddress, store)
				// if err != nil {
				// 	return err
				// }
			}
			i++
		}

		// for {
		// 	select {
		// 	case <-timerPool.C:
		// 		fmt.Println("WriteMetric")
		// 		mertic.WriteMetric(store)
		// 	case <-timerInterval.C:
		// 		fmt.Println("== SendMetric")
		// 		err := mertic.SendMetric(serverAddress, store)
		// 		if err != nil {
		// 			return err
		// 		}
		// 	}

	}
}
