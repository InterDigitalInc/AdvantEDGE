package util

import (
	"github.com/spf13/viper"
)

type Config struct {
	Mode          string `mapstructure: mode`
	SandboxUrl    string `mapstructure:"sandbox"`
	MecPlatform   string `mapstructure: "mecplatform"`
	SandboxName   string `mapstructure:"sandboxname"`
	AppInstanceId string `mapstructure:"appid"`
	Localurl      string `mapstructure:"localurl"`
	Port          string `mapstructure:"port"`
}

func LoadConfig(path string, name string) (config Config, err error) {
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return config, err
	}
	return

}
