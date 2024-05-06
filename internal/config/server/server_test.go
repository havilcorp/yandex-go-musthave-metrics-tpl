// Package server конфигурации сервера
package server

import "testing"

func TestConfig_WriteByFlag(t *testing.T) {
	conf := NewServerConfig()
	if err := conf.WriteByFlag(); (err != nil) != false {
		t.Errorf("Config.WriteByFlag() error = %v, wantErr %v", err, false)
	}
}

func TestConfig_WriteByEnv(t *testing.T) {
	conf := NewServerConfig()
	t.Run("Good", func(t *testing.T) {
		t.Setenv("ADDRESS", "ADDRESS")
		t.Setenv("STORE_INTERVAL", "300")
		t.Setenv("FILE_STORAGE_PATH", "/tmp/metrics-db.json")
		t.Setenv("RESTORE", "false")
		t.Setenv("DATABASE_DSN", "DATABASE_DSN")
		t.Setenv("KEY", "KEY")
		t.Setenv("CRYPTO_KEY", "CRYPTO_KEY")
		t.Setenv("CONFIG", "../../../config/server.json")
		if err := conf.WriteByEnv(); (err != nil) != false {
			t.Errorf("Config.WriteServerConfig() error = %v, wantErr %v", err, false)
		}
	})
	t.Run("Error STORE_INTERVAL", func(t *testing.T) {
		t.Setenv("ADDRESS", "ADDRESS")
		t.Setenv("STORE_INTERVAL", "not int")
		t.Setenv("FILE_STORAGE_PATH", "/tmp/metrics-db.json")
		t.Setenv("RESTORE", "false")
		t.Setenv("DATABASE_DSN", "DATABASE_DSN")
		t.Setenv("KEY", "KEY")
		t.Setenv("CRYPTO_KEY", "CRYPTO_KEY")
		t.Setenv("CONFIG", "../../../config/server.json")
		err := conf.WriteByEnv()
		if err == nil {
			t.Error("Config.WriteByEnv() not error")
		}
	})
}
