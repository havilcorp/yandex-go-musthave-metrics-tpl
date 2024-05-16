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
	AddressGRPC    string `json:"address_grpc"`
	CryptoKey      string `json:"crypto_key"`
	CryptoCrt      string `json:"crypto_crt"`
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
//   - -address_grpc - адрес grpc
//   - -r - интервал отправки метрик на сервер
//   - -p - интервал сбора метрик
//   - -k - ключ sha256
//   - -l - лимит запросов
//   - -crypto-key - путь к файлу с публичным ключем для шифрования сообщения
//   - -crypto-crt - путь к файлу с сертификатом
//   - -c - путь к файлу конфигов
func (c *Config) WriteByFlag() {
	flag.StringVar(&c.ServerAddress, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&c.AddressGRPC, "address_grpc", "", "address and port to run grpc server")
	flag.IntVar(&c.ReportInterval, "r", 10, "report interval time in sec")
	flag.IntVar(&c.PollInterval, "p", 2, "poll interval time in sec")
	flag.StringVar(&c.Key, "k", "", "sha256 key")
	flag.IntVar(&c.RateLimit, "l", 2, "rate limit")
	flag.StringVar(&c.CryptoKey, "crypto-key", "", "public key path")
	flag.StringVar(&c.CryptoCrt, "crypto-crt", "", "certificate path")
	flag.StringVar(&c.Config, "c", "", "config path")
	flag.Parse()
}

// WriteByEnv чтение настроек агента, env перекрывают флаги
//
// Env
//   - ADDRESS - адрес и порт сервера
//   - ADDRESS_GRPC - адрес и порт GRPC сервера
//   - REPORT_INTERVAL - интервал отправки метрик на сервер
//   - POLL_INTERVAL - интервал сбора метрик
//   - KEY - ключ sha256
//   - RATE_LIMIT - лимит запросов
//   - CRYPTO_KEY - путь до файла с публичным ключем для шифрования сообщения
//   - CRYPTO_CRT - путь к файлу с сертификатом
//   - CONFIG - путь к файлу конфигов
func (c *Config) WriteByEnv() error {
	if envServerAddress := os.Getenv("ADDRESS"); envServerAddress != "" {
		c.ServerAddress = envServerAddress
	}

	if envAddressGRPC := os.Getenv("ADDRESS_GRPC"); envAddressGRPC != "" {
		c.AddressGRPC = envAddressGRPC
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

	if envCryptoCrt := os.Getenv("CRYPTO_CRT"); envCryptoCrt != "" {
		c.CryptoCrt = envCryptoCrt
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
			AddressGRPC:    "",
			ReportInterval: 10,
			PollInterval:   2,
			RateLimit:      2,
			CryptoKey:      "",
			CryptoCrt:      "",
			Key:            "",
		}
		err = json.Unmarshal(data, &conf)
		if err != nil {
			return err
		}
		if c.ServerAddress == "localhost:8080" {
			c.ServerAddress = conf.ServerAddress
		}
		if c.AddressGRPC == "" {
			c.AddressGRPC = conf.AddressGRPC
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
		if c.CryptoCrt == "" {
			c.CryptoCrt = conf.CryptoCrt
		}
	}

	return nil
}
