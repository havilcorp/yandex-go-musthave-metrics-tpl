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

	if envAddress := os.Getenv("ADDRESS"); envAddress != "" {
		*flagServerAddr = envAddress
	}

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

func WriteServerConfig(serverAddress *string, storeInterval *int, fileStoragePath *string, isRestore *bool) {
	if envServerAddress := os.Getenv("ADDRESS"); envServerAddress != "" {
		*serverAddress = envServerAddress
	} else {
		flag.StringVar(serverAddress, "a", "localhost:8080", "address and port to run server")
	}

	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		envStoreIntervalVal, err := strconv.Atoi(envStoreInterval)
		if err != nil {
			panic(err)
		}
		*storeInterval = envStoreIntervalVal
	} else {
		flag.IntVar(storeInterval, "i", 300, "store save interval time in sec")
	}

	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		*fileStoragePath = envFileStoragePath
	} else {
		flag.StringVar(fileStoragePath, "f", "/tmp/metrics-db.json", "file store path save")
	}

	if envIsRestore := os.Getenv("RESTORE"); envIsRestore != "" {
		*isRestore = envIsRestore == "true"
	} else {
		flag.BoolVar(isRestore, "r", true, "is restore")
	}
	flag.Parse()
}
