/*
 * Copyright (c) 2022  The AdvantEDGE Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package util

import (
	"github.com/spf13/viper"
)

type Config struct {
	Mode          string `mapstructure:"mode"`
	SandboxUrl    string `mapstructure:"sandbox"`
	SandboxName   string `mapstructure:"sandboxname"`
	NodeName      string `mapstructure:"nodename"`
	HttpsOnly     bool   `mapstructure:"https"`
	MecPlatform   string `mapstructure:"mecplatform"`
	AppInstanceId string `mapstructure:"appid"`
	Localurl      string `mapstructure:"localurl"`
	Port          string `mapstructure:"port"`
	AppName       string `mapstructure:"appname"`
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
