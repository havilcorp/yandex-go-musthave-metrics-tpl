package config

import (
	"flag"
	"os"
	"strconv"
)

type Config struct {
	ServerAddress   string
	ReportInterval  int
	PollInterval    int
	StoreInterval   int
	FileStoragePath string
	IsRestore       bool
	DbConnect       string
}

func (c *Config) WriteAgentConfig() error {
	if envServerAddress := os.Getenv("ADDRESS"); envServerAddress != "" {
		c.ServerAddress = envServerAddress
	} else {
		flag.StringVar(&c.ServerAddress, "a", "localhost:8080", "address and port to run server")
	}

	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		envReportIntervalVal, err := strconv.Atoi(envReportInterval)
		if err != nil {
			return err
		}
		c.ReportInterval = envReportIntervalVal
	} else {
		flag.IntVar(&c.ReportInterval, "r", 10, "report interval time in sec")
	}

	if envPoolInterval := os.Getenv("POLL_INTERVAL"); envPoolInterval != "" {
		envPoolIntervalVal, err := strconv.Atoi(envPoolInterval)
		if err != nil {
			return err
		}
		c.PollInterval = envPoolIntervalVal
	} else {
		flag.IntVar(&c.PollInterval, "p", 2, "poll interval time in sec")
	}

	flag.Parse()

	return nil
}

func (c *Config) WriteServerConfig() error {

	if envServerAddress := os.Getenv("ADDRESS"); envServerAddress != "" {
		c.ServerAddress = envServerAddress
	} else {
		flag.StringVar(&c.ServerAddress, "a", "localhost:8080", "address and port to run server")
	}

	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		envStoreIntervalVal, err := strconv.Atoi(envStoreInterval)
		if err != nil {
			return err
		}
		c.StoreInterval = envStoreIntervalVal
	} else {
		flag.IntVar(&c.StoreInterval, "i", 300, "store save interval time in sec")
	}

	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		c.FileStoragePath = envFileStoragePath
	} else {
		flag.StringVar(&c.FileStoragePath, "f", "/tmp/metrics-db.json", "file store path save")
	}

	if envIsRestore := os.Getenv("RESTORE"); envIsRestore != "" {
		c.IsRestore = envIsRestore == "true"
	} else {
		flag.BoolVar(&c.IsRestore, "r", true, "is restore")
	}

	if envDbConnect := os.Getenv("DATABASE_DSN"); envDbConnect != "" {
		c.DbConnect = envDbConnect
	} else {
		flag.StringVar(
			&c.DbConnect,
			"d",
			"", // postgres://postgres:password@localhost:5433/postgres?sslmode=disable
			"db connect string",
		)
	}

	flag.Parse()

	return nil
}
