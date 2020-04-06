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
	"time"

	"github.com/InterDigitalInc/AdvantEDGE/go-apps/meepctl/utils"

	"github.com/roymx/viper"
	"github.com/spf13/cobra"
)

// configSet represents the set command
var configPorts = &cobra.Command{
	Use:     "ports [http port #/https port #]",
	Short:   "get/set HTTP/HTTPS ports in the meepctl config file",
	Long:    "Get/Set HTTP/HTTPS ports in the meepctl config file",
	Example: "  meepctl config ports 32080/32443",
	Args:    cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		key := "meep.ports"
		v, _ := cmd.Flags().GetBool("verbose")
		t, _ := cmd.Flags().GetBool("time")
		if v {
			fmt.Println("config ports called")
			fmt.Println("[flag] verbose:", v)
			fmt.Println("[flag] time:", t)
		}

		start := time.Now()
		value := viper.GetString(key)
		if len(args) == 0 {
			_ = cmd.Help()
			fmt.Println("")
		} else {
			ports := args[0]
			valid, reason := utils.ConfigPortsValid(ports)
			if valid {
				utils.Cfg.Meep.Ports = ports
				err := utils.ConfigWriteFile(utils.Cfg, viper.ConfigFileUsed())
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("Updated meep.ports with [" + ports + "]")
				value = ports
			} else {
				fmt.Println("Invalid Ports: " + reason)
				fmt.Println("")
				_ = cmd.Help()
			}
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
	configCmd.AddCommand(configPorts)
}
