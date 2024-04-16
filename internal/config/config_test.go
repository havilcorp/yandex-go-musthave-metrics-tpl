package config

import "fmt"

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
