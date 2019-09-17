/*
 * Copyright (c) 2019  InterDigital Communications, Inc
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

package cmd

import (
	"fmt"
	"sort"

	"github.com/roymx/viper"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "manage meep environment configuration",
	Long: `Get/Set meep environment configuration.
Config file location: ~/.meepctl.yaml

Config file is created with default values if it does not already exist.
Values may be changed using the available commands described below.`,

	Run: func(cmd *cobra.Command, args []string) {
		keys := viper.AllKeys()
		sort.Strings(keys)

		_ = cmd.Help()
		fmt.Println("")
		fmt.Println("========================================")
		for _, key := range keys {
			fmt.Println(key, ":", viper.GetString(key))
		}
		fmt.Println("========================================")
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
