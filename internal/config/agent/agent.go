// Package agent конфигурации агента
package agent

import (
	"encoding/json"
	"flag"
	"os"
	"strconv"
)

type Config struct {
	Config         string
	ServerAddress  string `json:"address"`
	CryptoKey      string `json:"crypto_key"`
	Key            string `json:"key"`
	ReportInterval int    `json:"report_interval"`
	PollInterval   int    `json:"poll_interval"`
	RateLimit      int    `json:"rate_limit"`
}

func NewAgentConfig() *Config {
	return &Config{}
}

// WriteByFlag чтение настроек агента через флаги
//
// Флаги
//   - -a - адрес и порт сервера
//   - -r - интервал отправки метрик на сервер
//   - -p - интервал сбора метрик
//   - -k - ключ sha256
//   - -l - лимит запросов
//   - -crypto-key - путь к файлу с публичным ключем для шифрования сообщения
//   - -c - путь к файлу конфигов
func (c *Config) WriteByFlag() error {
	flag.StringVar(&c.ServerAddress, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&c.ReportInterval, "r", 10, "report interval time in sec")
	flag.IntVar(&c.PollInterval, "p", 2, "poll interval time in sec")
	flag.StringVar(&c.Key, "k", "", "sha256 key")
	flag.IntVar(&c.RateLimit, "l", 2, "rate limit")
	flag.StringVar(&c.CryptoKey, "crypto-key", "", "public key path")
	flag.StringVar(&c.Config, "c", "", "config path")
	flag.Parse()
	return nil
}

// WriteByEnv чтение настроек агента, env перекрывают флаги
//
// Env
//   - ADDRESS - адрес и порт сервера
//   - REPORT_INTERVAL - интервал отправки метрик на сервер
//   - POLL_INTERVAL - интервал сбора метрик
//   - KEY - ключ sha256
//   - RATE_LIMIT - лимит запросов
//   - CRYPTO_KEY - путь до файла с публичным ключем для шифрования сообщения
//   - CONFIG - путь к файлу конфигов
func (c *Config) WriteByEnv() error {
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
			return err
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

	if envConfig := os.Getenv("CONFIG"); envConfig != "" {
		c.Config = envConfig
	}

	if c.Config != "" {
		data, err := os.ReadFile(c.Config)
		if err != nil {
			return err
		}
		conf := Config{
			ServerAddress:  "localhost:8080",
			ReportInterval: 10,
			PollInterval:   2,
			RateLimit:      2,
			CryptoKey:      "",
			Key:            "",
		}
		err = json.Unmarshal(data, &conf)
		if err != nil {
			return err
		}
		if c.ServerAddress == "localhost:8080" {
			c.ServerAddress = conf.ServerAddress
		}
		if c.Key == "" {
			c.Key = conf.Key
		}
		if c.ReportInterval == 10 {
			c.ReportInterval = conf.ReportInterval
		}
		if c.PollInterval == 2 {
			c.PollInterval = conf.PollInterval
		}
		if c.RateLimit == 2 {
			c.RateLimit = conf.RateLimit
		}
		if c.CryptoKey == "" {
			c.CryptoKey = conf.CryptoKey
		}
	}

	return nil
}
