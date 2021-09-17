package util

import (
	"github.com/spf13/viper"
)

// Config stores all configuration of Demo 3
// The values are read by viper from a config file or env file
type Config struct {
	SandboxUrl    string `mapstructure:"sandbox"`
	AppInstanceId string `mapstructure:"appid"`
	Port          string `mapstructure:"port"`
	ServiceName   string `mapstructure:"service"`
}

// LoadConfig reads configuration from a environment variable specified by path
func LoadConfig(path string, name string) (config Config, err error) {
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	viper.SetConfigName(name)

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
