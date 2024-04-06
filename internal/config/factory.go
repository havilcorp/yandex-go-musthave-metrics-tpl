package config

// func ConfigFactory(provider string) (*Config, error) {
// 	conf := NewConfig()
// 	switch provider {
// 	case "agent":
// 		if err := conf.WriteAddressConfig(); err != nil {
// 			logrus.Error(err)
// 			return nil, err
// 		}
// 		if err := conf.WriteSHA256Config(); err != nil {
// 			logrus.Error(err)
// 			return nil, err
// 		}
// 		if err := conf.WriteAgentConfig(); err != nil {
// 			logrus.Error(err)
// 			return nil, err
// 		}
// 		return conf, nil
// 	case "server":
// 		if err := conf.WriteAddressConfig(); err != nil {
// 			logrus.Error(err)
// 			return nil, err
// 		}
// 		if err := conf.WriteSHA256Config(); err != nil {
// 			logrus.Error(err)
// 			return nil, err
// 		}
// 		if err := conf.WriteServerConfig(); err != nil {
// 			logrus.Error(err)
// 			return nil, err
// 		}
// 		if err := conf.WriteIsRestoreConfig(); err != nil {
// 			logrus.Error(err)
// 			return nil, err
// 		}
// 		return conf, nil
// 	default:
// 		return nil, fmt.Errorf("unknown provider %s", provider)
// 	}
// }
