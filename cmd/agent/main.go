package main

import (
	"fmt"
	"runtime"
	"time"
)

type Metric struct {
	gauge   float64
	counter int64
}

type MemStorage struct {
	list map[string]Metric
}

func main() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	fmt.Println("Alloc", memStats.Alloc)
	fmt.Println("BuckHashSys", memStats.BuckHashSys)
	fmt.Println("Frees", memStats.Frees)
	fmt.Println("GCCPUFraction", memStats.GCCPUFraction)
	fmt.Println("GCSys", memStats.GCSys)
	fmt.Println("HeapAlloc", memStats.HeapAlloc)
	fmt.Println("HeapIdle", memStats.HeapIdle)
	fmt.Println("HeapInuse", memStats.HeapInuse)
	fmt.Println("HeapObjects", memStats.HeapObjects)
	fmt.Println("HeapReleased", memStats.HeapReleased)
	fmt.Println("HeapSys", memStats.HeapSys)
	fmt.Println("LastGC", memStats.LastGC)
	fmt.Println("Lookups", memStats.Lookups)
	fmt.Println("MCacheInuse", memStats.MCacheInuse)
	fmt.Println("MCacheSys", memStats.MCacheSys)
	fmt.Println("MSpanInuse", memStats.MSpanInuse)
	fmt.Println("MSpanSys", memStats.MSpanSys)
	fmt.Println("Mallocs", memStats.Mallocs)
	fmt.Println("NextGC", memStats.NextGC)
	fmt.Println("NumForcedGC", memStats.NumForcedGC)
	fmt.Println("NumGC", memStats.NumGC)
	fmt.Println("OtherSys", memStats.OtherSys)
	fmt.Println("PauseTotalNs", memStats.PauseTotalNs)
	fmt.Println("StackInuse", memStats.StackInuse)
	fmt.Println("StackSys", memStats.StackSys)
	fmt.Println("Sys", memStats.Sys)
	fmt.Println("TotalAlloc", memStats.TotalAlloc)

	a := MemStorage{list: map[string]Metric{}}
	a.list["Alloc"] = Metric{gauge: 1, counter: 2}
	fmt.Println(a)

	c := time.Tick(2 * time.Second)
	for next := range c {
		fmt.Printf("%v\n", next)
	}
}
