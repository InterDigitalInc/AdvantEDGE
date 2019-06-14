/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
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
var configIp = &cobra.Command{
	Use:   "ip [IP]",
	Short: "get/get node IP address in the meepctl config file",
	Long:  `Get/Set node IP address in the meepctl config file`,
	Example: `  # Get currently configured node IP address
  meepctl config ip
  # Set node IP Address to 1.2.3.4
  meepctl config ip 1.2.3.4`,
	Run: func(cmd *cobra.Command, args []string) {
		v, _ := cmd.Flags().GetBool("verbose")
		t, _ := cmd.Flags().GetBool("time")
		if v {
			fmt.Println("config ip called")
			fmt.Println("[flag] verbose:", v)
			fmt.Println("[flag] time:", t)
		}

		start := time.Now()
		if len(args) > 0 {
			ip := args[0]
			valid, reason := utils.ConfigIPValid(ip)
			if valid {
				cfg := utils.ConfigReadFile(viper.ConfigFileUsed())
				cfg.Node.IP = ip
				err := utils.ConfigWriteFile(cfg, viper.ConfigFileUsed())
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Println("Updated node.ip with [" + ip + "]")
			} else {
				fmt.Println("Invalid IP: " + reason)
				fmt.Println("")
				_ = cmd.Help()
			}
		} else {
			key := "node.ip"
			fmt.Println(key, ":", viper.GetString(key))
		}

		elapsed := time.Since(start)
		if t {
			fmt.Println("Took ", elapsed.Round(time.Millisecond).String())
		}
	},
}

func init() {
	configCmd.AddCommand(configIp)
}
