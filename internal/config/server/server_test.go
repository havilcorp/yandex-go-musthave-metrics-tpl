// Package server конфигурации сервера
package server

import "testing"

func TestConfig_WriteServerConfig(t *testing.T) {
	conf := NewServerConfig()
	t.Run("Test server", func(t *testing.T) {
		t.Setenv("ADDRESS", "ADDRESS")
		t.Setenv("STORE_INTERVAL", "300")
		t.Setenv("FILE_STORAGE_PATH", "/tmp/metrics-db.json")
		t.Setenv("RESTORE", "false")
		t.Setenv("DATABASE_DSN", "DATABASE_DSN")
		t.Setenv("KEY", "KEY")
		t.Setenv("CRYPTO_KEY", "CRYPTO_KEY")
		if err := conf.WriteServerConfig(); (err != nil) != false {
			t.Errorf("Config.WriteServerConfig() error = %v, wantErr %v", err, false)
		}
	})
}
