// Package agent конфигурации агента
package agent

import (
	"testing"
)

func TestConfig_WriteByFlag(t *testing.T) {
	conf := NewAgentConfig()
	if err := conf.WriteByFlag(); (err != nil) != false {
		t.Errorf("Config.WriteByFlag() error = %v, wantErr %v", err, false)
	}
}

func TestConfig_WriteByEnv(t *testing.T) {
	conf := NewAgentConfig()
	t.Run("Good", func(t *testing.T) {
		t.Setenv("ADDRESS", "ADDRESS")
		t.Setenv("REPORT_INTERVAL", "10")
		t.Setenv("POLL_INTERVAL", "2")
		t.Setenv("KEY", "KEY")
		t.Setenv("RATE_LIMIT", "2")
		t.Setenv("CRYPTO_KEY", "CRYPTO_KEY")
		t.Setenv("CONFIG", "../../../config/agent.json")
		if err := conf.WriteByEnv(); (err != nil) != false {
			t.Errorf("Config.WriteByEnv() error = %v, wantErr %v", err, false)
		}
	})
	t.Run("Error REPORT_INTERVAL", func(t *testing.T) {
		t.Setenv("ADDRESS", "ADDRESS")
		t.Setenv("REPORT_INTERVAL", "not int")
		t.Setenv("POLL_INTERVAL", "2")
		t.Setenv("KEY", "KEY")
		t.Setenv("RATE_LIMIT", "2")
		t.Setenv("CRYPTO_KEY", "CRYPTO_KEY")
		t.Setenv("CONFIG", "../../../config/agent.json")
		err := conf.WriteByEnv()
		if err == nil {
			t.Error("Config.WriteByEnv() not error")
		}
	})
	t.Run("Error POLL_INTERVAL", func(t *testing.T) {
		t.Setenv("ADDRESS", "ADDRESS")
		t.Setenv("REPORT_INTERVAL", "10")
		t.Setenv("POLL_INTERVAL", "not int")
		t.Setenv("KEY", "KEY")
		t.Setenv("RATE_LIMIT", "2")
		t.Setenv("CRYPTO_KEY", "CRYPTO_KEY")
		t.Setenv("CONFIG", "../../../config/agent.json")
		err := conf.WriteByEnv()
		if err == nil {
			t.Error("Config.WriteByEnv() not error")
		}
	})
	t.Run("Error RATE_LIMIT", func(t *testing.T) {
		t.Setenv("ADDRESS", "ADDRESS")
		t.Setenv("REPORT_INTERVAL", "10")
		t.Setenv("POLL_INTERVAL", "2")
		t.Setenv("KEY", "KEY")
		t.Setenv("RATE_LIMIT", "not int")
		t.Setenv("CRYPTO_KEY", "CRYPTO_KEY")
		t.Setenv("CONFIG", "../../../config/agent.json")
		err := conf.WriteByEnv()
		if err == nil {
			t.Error("Config.WriteByEnv() not error")
		}
	})
}
