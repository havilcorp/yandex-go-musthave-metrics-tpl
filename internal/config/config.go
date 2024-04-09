// Конфигурации агента и сервера
package config

import (
	"flag"
	"fmt"
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
	DBConnect       string
	Key             string
	RateLimit       int
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) String() string {
	out := ""
	out += "\n******* CONFIG *******\n"
	out += fmt.Sprintf("* ServerAddress: %s\n", c.ServerAddress)
	out += fmt.Sprintf("* ReportInterval: %d\n", c.ReportInterval)
	out += fmt.Sprintf("* PollInterval: %d\n", c.PollInterval)
	out += fmt.Sprintf("* StoreInterval: %d\n", c.StoreInterval)
	out += fmt.Sprintf("* FileStoragePath: %s\n", c.FileStoragePath)
	out += fmt.Sprintf("* IsRestore: %t\n", c.IsRestore)
	out += fmt.Sprintf("* DBConnect: %s\n", c.DBConnect)
	out += fmt.Sprintf("* Key: %s\n", c.Key)
	out += fmt.Sprintf("* RateLimit: %d", c.RateLimit)
	out += "\n**********************\n"
	return out
}

// WriteAgentConfig чтение настроек агента, env перекрывают флаги
//
// Флаги
//   - -a - адрес и порт сервера
//   - -r - интервал отправки метрик на сервер
//   - -p - интервал сбора метрик
//   - -k - ключ sha256
//   - -l - лимит запросов
//
// Env
//   - ADDRESS - адрес и порт сервера
//   - REPORT_INTERVAL - интервал отправки метрик на сервер
//   - POLL_INTERVAL - интервал сбора метрик
//   - KEY - ключ sha256
//   - RATE_LIMIT - лимит запросов
func (c *Config) WriteAgentConfig() error {
	flag.StringVar(&c.ServerAddress, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&c.ReportInterval, "r", 10, "report interval time in sec")
	flag.IntVar(&c.PollInterval, "p", 2, "poll interval time in sec")
	flag.StringVar(&c.Key, "k", "", "sha256 key")
	flag.IntVar(&c.RateLimit, "l", 2, "rate limit")
	flag.Parse()

	if envServerAddress := os.Getenv("ADDRESS"); envServerAddress != "" {
		c.ServerAddress = envServerAddress
	}

	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		envReportIntervalVal, err := strconv.Atoi(envReportInterval)
		if err != nil {
			return err
		}
		c.ReportInterval = envReportIntervalVal
	}

	if envPoolInterval := os.Getenv("POLL_INTERVAL"); envPoolInterval != "" {
		envPoolIntervalVal, err := strconv.Atoi(envPoolInterval)
		if err != nil {
			return nil
		}
		c.PollInterval = envPoolIntervalVal
	}

	if envKey := os.Getenv("KEY"); envKey != "" {
		c.Key = envKey
	}

	if envRateLimit := os.Getenv("RATE_LIMIT"); envRateLimit != "" {
		envRateLimitVal, err := strconv.Atoi(envRateLimit)
		if err != nil {
			return err
		}
		c.RateLimit = envRateLimitVal
	}

	return nil
}

// WriteAgentConfig чтение настроек сервера, env перекрывают флаги
//
// Флаги
//   - -a - адрес и порт сервера
//   - -i - интевал сохранения метрик в файл
//   - -f - файл для созранения метрик. По деволту: /tmp/metrics-db.json
//   - -r - загружать ли при запуске метрики из файла
//   - -d - строка подключения к базе данных
//   - -r - ключ sha256
//
// Env
//   - ADDRESS - адрес и порт сервера
//   - STORE_INTERVAL - интевал сохранения метрик в файл
//   - FILE_STORAGE_PATH - файл для созранения метрик. По деволту: /tmp/metrics-db.json
//   - RESTORE - загружать ли при запуске метрики из файла
//   - DATABASE_DSN - строка подключения к базе данных
//   - KEY - ключ sha256
func (c *Config) WriteServerConfig() error {
	flag.StringVar(&c.ServerAddress, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&c.StoreInterval, "i", 300, "store save interval time in sec")
	flag.StringVar(&c.FileStoragePath, "f", "/tmp/metrics-db.json", "file store path save")
	flag.BoolVar(&c.IsRestore, "r", true, "is restore")
	flag.StringVar(&c.DBConnect, "d", "", "db connect string")
	flag.StringVar(&c.Key, "k", "", "sha256 key")
	flag.Parse()

	if envServerAddress := os.Getenv("ADDRESS"); envServerAddress != "" {
		c.ServerAddress = envServerAddress
	}

	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		envStoreIntervalVal, err := strconv.Atoi(envStoreInterval)
		if err != nil {
			return err
		}
		c.StoreInterval = envStoreIntervalVal
	}

	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		c.FileStoragePath = envFileStoragePath
	}

	if envIsRestore := os.Getenv("RESTORE"); envIsRestore != "" {
		c.IsRestore = envIsRestore == "true"
	}

	if envDBConnect := os.Getenv("DATABASE_DSN"); envDBConnect != "" {
		c.DBConnect = envDBConnect
	}

	if envKey := os.Getenv("KEY"); envKey != "" {
		c.Key = envKey
	}

	return nil
}
