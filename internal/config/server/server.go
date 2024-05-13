// Package server конфигурации сервера
package server

import (
	"encoding/json"
	"flag"
	"os"
	"strconv"
)

type Config struct {
	Config          string
	ServerAddress   string `json:"address"`
	FileStoragePath string `json:"store_file"`
	DBConnect       string `json:"database_dsn"`
	Key             string `json:"key"`
	CryptoKey       string `json:"crypto_key"`
	IsRestore       bool   `json:"restore"`
	StoreInterval   int    `json:"store_interval"`
	TrustedSubnet   string `json:"trusted_subnet"`
}

func NewServerConfig() *Config {
	return &Config{}
}

// WriteByFlag чтение настроек сервера через флаги
//
// Флаги
//   - -a - адрес и порт сервера
//   - -i - интевал сохранения метрик в файл
//   - -f - файл для созранения метрик. По деволту: /tmp/metrics-db.json
//   - -r - загружать ли при запуске метрики из файла
//   - -d - строка подключения к базе данных
//   - -r - ключ sha256
//   - -crypto-key - путь к файлу с приватным ключем для расшифрования сообщения
//   - -c - путь к файлу конфигов
//   - -t - доверенная маска подсети
func (c *Config) WriteByFlag() error {
	flag.StringVar(&c.ServerAddress, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&c.StoreInterval, "i", 300, "store save interval time in sec")
	flag.StringVar(&c.FileStoragePath, "f", "/tmp/metrics-db.json", "file store path save")
	flag.BoolVar(&c.IsRestore, "r", true, "is restore")
	flag.StringVar(&c.DBConnect, "d", "", "db connect string")
	flag.StringVar(&c.Key, "k", "", "sha256 key")
	flag.StringVar(&c.CryptoKey, "crypto-key", "", "private key path")
	flag.StringVar(&c.Config, "c", "", "config path")
	flag.StringVar(&c.TrustedSubnet, "t", "", "trusted subnet")
	flag.Parse()
	return nil
}

// WriteByEnv чтение настроек сервера, env перекрывают флаги
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
//   - TRUSTED_SUBNET - доверенная маска подсети
func (c *Config) WriteByEnv() error {
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

	if envConfig := os.Getenv("CONFIG"); envConfig != "" {
		c.Config = envConfig
	}

	if envTrustedSubnet := os.Getenv("TRUSTED_SUBNET"); envTrustedSubnet != "" {
		c.TrustedSubnet = envTrustedSubnet
	}

	if c.Config != "" {
		data, err := os.ReadFile(c.Config)
		if err != nil {
			return err
		}
		conf := Config{
			ServerAddress:   "localhost:8080",
			StoreInterval:   300,
			FileStoragePath: "/tmp/metrics-db.json",
			IsRestore:       true,
			DBConnect:       "",
			CryptoKey:       "",
			Key:             "",
			TrustedSubnet:   "",
		}
		err = json.Unmarshal(data, &conf)
		if err != nil {
			return err
		}
		if c.ServerAddress == "localhost:8080" {
			c.ServerAddress = conf.ServerAddress
		}
		if c.StoreInterval == 300 {
			c.StoreInterval = conf.StoreInterval
		}
		if c.FileStoragePath == "" {
			c.FileStoragePath = conf.FileStoragePath
		}
		if c.IsRestore {
			c.IsRestore = conf.IsRestore
		}
		if c.DBConnect == "" {
			c.DBConnect = conf.DBConnect
		}
		if c.Key == "" {
			c.Key = conf.Key
		}
		if c.CryptoKey == "" {
			c.CryptoKey = conf.CryptoKey
		}
		if c.TrustedSubnet == "" {
			c.TrustedSubnet = conf.TrustedSubnet
		}
	}

	return nil
}
