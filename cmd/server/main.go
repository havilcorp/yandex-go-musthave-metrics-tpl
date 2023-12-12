package main

import (
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/server"
)

func main() {
	err := server.StartServer()
	if err != nil {
		panic(err)
	}
}
