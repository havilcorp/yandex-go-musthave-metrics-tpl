package main

import (
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/agent"
)

func main() {

	err := agent.StartAgent()
	if err != nil {
		panic(err)
	}

}
