package config

import (
	"flag"
	"os"
	"strconv"
)

func WriteAgentConfig(flagServerAddr *string, reportInterval *int, pollInterval *int) {
	flag.StringVar(flagServerAddr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(reportInterval, "r", 10, "report interval time in sec")
	flag.IntVar(pollInterval, "p", 2, "poll interval time in sec")

	flag.Parse()

	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		envReportIntervalVal, err := strconv.Atoi(envReportInterval)
		if err != nil {
			panic(err)
		}
		*reportInterval = envReportIntervalVal
	}

	if envPoolInterval := os.Getenv("POLL_INTERVAL"); envPoolInterval != "" {
		envPoolIntervalVal, err := strconv.Atoi(envPoolInterval)
		if err != nil {
			panic(err)
		}
		*pollInterval = envPoolIntervalVal
	}

}

func WriteServerConfig(serverAddress *string) {
	flag.StringVar(serverAddress, "a", "localhost:8080", "address and port to run server")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		*serverAddress = envRunAddr
	}
}
