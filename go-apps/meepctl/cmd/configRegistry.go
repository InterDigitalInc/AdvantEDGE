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

package cmd

import (
	"fmt"
	"time"

	"github.com/InterDigitalInc/AdvantEDGE/go-apps/meepctl/utils"

	"github.com/roymx/viper"
	"github.com/spf13/cobra"
)

// configSet represents the set command
var configRegistry = &cobra.Command{
	Use:     "registry [name:port]",
	Short:   "get/set docker registry meepctl config file",
	Long:    "Get/Set docker registry in the meepctl config file",
	Example: "  meepctl config registry meep-docker-registry:30001",
	Args:    cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		key := "meep.registry"
		v, _ := cmd.Flags().GetBool("verbose")
		t, _ := cmd.Flags().GetBool("time")
		if v {
			fmt.Println("config registry called")
			fmt.Println("[flag] verbose:", v)
			fmt.Println("[flag] time:", t)
		}

		start := time.Now()
		value := viper.GetString(key)

		if len(args) == 0 {
			_ = cmd.Help()
			fmt.Println("")
		} else {
			registry := args[0]
			utils.Cfg.Meep.Registry = registry
			err := utils.ConfigWriteFile(utils.Cfg, viper.ConfigFileUsed())
			if err != nil {
				fmt.Println(err)
			}
			value = registry
		}
		fmt.Println("========================================")
		fmt.Println(key, ":", value)
		fmt.Println("========================================")

		elapsed := time.Since(start)
		if t {
			fmt.Println("Took ", elapsed.Round(time.Millisecond).String())
		}
	},
}

func init() {
	configCmd.AddCommand(configRegistry)
}
