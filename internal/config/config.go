// Package config конфигурации агента и сервера
package config

import (
	"encoding/json"
	"flag"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

type Config struct {
	ServerAddress   string `json:"address"`
	Key             string
	FileStoragePath string `json:"store_file"`
	DBConnect       string `json:"database_dsn"`
	ReportInterval  int    `json:"report_interval"`
	PollInterval    int    `json:"poll_interval"`
	StoreInterval   int    `json:"store_interval"`
	IsRestore       bool   `json:"restore"`
	RateLimit       int
	CryptoKey       string `json:"crypto_key"`
	Config          string
}

func NewConfig() *Config {
	return &Config{}
}

// WriteAgentConfig чтение настроек агента, env перекрывают флаги
//
// Флаги
//   - -a - адрес и порт сервера
//   - -r - интервал отправки метрик на сервер
//   - -p - интервал сбора метрик
//   - -k - ключ sha256
//   - -l - лимит запросов
//   - -crypto-key - путь к файлу с публичным ключем для шифрования сообщения
//   - -config - путь к файлу конфигов
//
// Env
//   - ADDRESS - адрес и порт сервера
//   - REPORT_INTERVAL - интервал отправки метрик на сервер
//   - POLL_INTERVAL - интервал сбора метрик
//   - KEY - ключ sha256
//   - RATE_LIMIT - лимит запросов
//   - CRYPTO_KEY - путь до файла с публичным ключем для шифрования сообщения
//   - CONFIG - путь к файлу конфигов
func (c *Config) WriteAgentConfig() error {
	flag.StringVar(&c.ServerAddress, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&c.ReportInterval, "r", 10, "report interval time in sec")
	flag.IntVar(&c.PollInterval, "p", 2, "poll interval time in sec")
	flag.StringVar(&c.Key, "k", "", "sha256 key")
	flag.IntVar(&c.RateLimit, "l", 2, "rate limit")
	flag.StringVar(&c.CryptoKey, "crypto-key", "", "public key path")
	flag.StringVar(&c.Config, "config", "", "config path")
	flag.Parse()

	if c.Config != "" {
		data, err := os.ReadFile(c.Config)
		if err != nil {
			return err
		}
		logrus.Info(data)
		var conf Config
		err = json.Unmarshal(data, &conf)
		if err != nil {
			return err
		}
		logrus.Info(conf.ServerAddress)
	}

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

	if envCryptoKey := os.Getenv("CRYPTO_KEY"); envCryptoKey != "" {
		c.CryptoKey = envCryptoKey
	}

	return nil
}

// WriteServerConfig чтение настроек сервера, env перекрывают флаги
//
// Флаги
//   - -a - адрес и порт сервера
//   - -i - интевал сохранения метрик в файл
//   - -f - файл для созранения метрик. По деволту: /tmp/metrics-db.json
//   - -r - загружать ли при запуске метрики из файла
//   - -d - строка подключения к базе данных
//   - -r - ключ sha256
//   - -crypto-key - путь к файлу с приватным ключем для расшифрования сообщения
//   - -config - путь к файлу конфигов
//
// Env
//   - ADDRESS - адрес и порт сервера
//   - STORE_INTERVAL - интевал сохранения метрик в файл
//   - FILE_STORAGE_PATH - файл для созранения метрик. По деволту: /tmp/metrics-db.json
//   - RESTORE - загружать ли при запуске метрики из файла
//   - DATABASE_DSN - строка подключения к базе данных
//   - KEY - ключ sha256
//   - CRYPTO_KEY - путь до файла с приватным ключем для расшифрования сообщения
//   - CONFIG - путь к файлу конфигов
func (c *Config) WriteServerConfig() error {
	flag.StringVar(&c.ServerAddress, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&c.StoreInterval, "i", 300, "store save interval time in sec")
	flag.StringVar(&c.FileStoragePath, "f", "/tmp/metrics-db.json", "file store path save")
	flag.BoolVar(&c.IsRestore, "r", true, "is restore")
	flag.StringVar(&c.DBConnect, "d", "", "db connect string")
	flag.StringVar(&c.Key, "k", "", "sha256 key")
	flag.StringVar(&c.CryptoKey, "crypto-key", "", "private key path")
	flag.StringVar(&c.Config, "config", "", "config path")
	flag.Parse()

	if c.Config != "" {
		data, err := os.ReadFile(c.Config)
		if err != nil {
			return err
		}
		var conf Config
		err = json.Unmarshal(data, &conf)
		if err != nil {
			return err
		}
		logrus.Info(conf.ServerAddress)
	}

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

	if envCryptoKey := os.Getenv("CRYPTO_KEY"); envCryptoKey != "" {
		c.CryptoKey = envCryptoKey
	}

	return nil
}
