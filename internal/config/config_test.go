package config

import (
	"fmt"
	"testing"
)

func ExampleConfig_WriteAgentConfig() {
	conf := NewConfig()
	err := conf.WriteAgentConfig()
	if err != nil {
		fmt.Print(err)
		return
	}
}

func ExampleConfig_WriteServerConfig() {
	conf := NewConfig()
	err := conf.WriteServerConfig()
	if err != nil {
		fmt.Print(err)
		return
	}
}

func TestConfig_WriteServerConfig(t *testing.T) {
	c := NewConfig()
	t.Run("Test server", func(t *testing.T) {
		t.Setenv("ADDRESS", "ADDRESS")
		t.Setenv("STORE_INTERVAL", "10")
		t.Setenv("FILE_STORAGE_PATH", "/tmp/metrics-db.json")
		t.Setenv("RESTORE", "false")
		t.Setenv("DATABASE_DSN", "DATABASE_DSN")
		t.Setenv("KEY", "KEY")
		if err := c.WriteServerConfig(); (err != nil) != false {
			t.Errorf("Config.WriteServerConfig() error = %v, wantErr %v", err, false)
		}
	})
}
