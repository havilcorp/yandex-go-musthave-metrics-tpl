// Package agent конфигурации агента
package agent

import (
	"encoding/json"
	"flag"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

type Config struct {
	ServerAddress  string `json:"address"`
	ReportInterval int    `json:"report_interval"`
	PollInterval   int    `json:"poll_interval"`
	Key            string
	RateLimit      int
	CryptoKey      string `json:"crypto_key"`
	Config         string
}

func NewAgentConfig() *Config {
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
//   - -c - путь к файлу конфигов
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
	flag.StringVar(&c.Config, "c", "", "config path")
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

	if envCryptoKey := os.Getenv("CRYPTO_KEY"); envCryptoKey != "" {
		c.CryptoKey = envCryptoKey
	}

	if c.Config != "" {
		data, err := os.ReadFile(c.Config)
		if err != nil {
			return err
		}
		logrus.Info(string(data))
		conf := Config{
			ServerAddress:  "localhost:8080",
			ReportInterval: 10,
			PollInterval:   2,
			CryptoKey:      "",
		}
		err = json.Unmarshal(data, &conf)
		if err != nil {
			return err
		}
		if c.ServerAddress == "localhost:8080" {
			c.ServerAddress = conf.ServerAddress
		}
		if c.ReportInterval == 10 {
			c.ReportInterval = conf.ReportInterval
		}
		if c.PollInterval == 2 {
			c.PollInterval = conf.PollInterval
		}
		if c.CryptoKey == "" {
			c.CryptoKey = conf.CryptoKey
		}
	}

	return nil
}
