// Package agent конфигурации агента
package agent

import (
	"testing"
)

func TestConfig_WriteAgentConfig(t *testing.T) {
	conf := NewAgentConfig()
	t.Run("Test server", func(t *testing.T) {
		t.Setenv("ADDRESS", "ADDRESS")
		t.Setenv("REPORT_INTERVAL", "10")
		t.Setenv("POLL_INTERVAL", "2")
		t.Setenv("KEY", "KEY")
		t.Setenv("RATE_LIMIT", "2")
		t.Setenv("CRYPTO_KEY", "CRYPTO_KEY")
		t.Setenv("CONFIG", "../../../config/agent.json")
		if err := conf.WriteAgentConfig(); (err != nil) != false {
			t.Errorf("Config.WriteAgentConfig() error = %v, wantErr %v", err, false)
		}
	})
}
