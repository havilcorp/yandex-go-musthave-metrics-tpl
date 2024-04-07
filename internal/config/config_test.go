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

// func TestConfig_WriteAddressConfig(t *testing.T) {
// 	t.Run("WriteAddressConfig", func(t *testing.T) {
// 		conf := Config{}
// 		err := conf.WriteAddressConfig()
// 		if err != nil {
// 			t.Errorf("WriteAddressConfig() error = %v", err)
// 		}
// 		require.NotEmpty(t, conf.String())
// 	})
// }

// func TestConfig_WriteSHA256Config(t *testing.T) {
// 	t.Run("WriteSHA256Config", func(t *testing.T) {
// 		conf := Config{}
// 		err := conf.WriteSHA256Config()
// 		if err != nil {
// 			t.Errorf("WriteSHA256Config() error = %v", err)
// 		}
// 		require.NotEmpty(t, conf.String())
// 	})
// }

// func TestConfig_WriteAgentConfig(t *testing.T) {
// 	t.Run("WriteAgentConfig", func(t *testing.T) {
// 		conf := Config{}
// 		err := conf.WriteAgentConfig()
// 		if err != nil {
// 			t.Errorf("WriteAgentConfig() error = %v", err)
// 		}
// 		require.NotEmpty(t, conf.String())
// 	})
// }

// func TestConfig_WriteServerConfig(t *testing.T) {
// 	t.Run("WriteServerConfig", func(t *testing.T) {
// 		conf := Config{}
// 		err := conf.WriteServerConfig()
// 		if err != nil {
// 			t.Errorf("WriteServerConfig() error = %v", err)
// 		}
// 		require.NotEmpty(t, conf.String())
// 	})
// }
