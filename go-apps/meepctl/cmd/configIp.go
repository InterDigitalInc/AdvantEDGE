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
	Use:     "ip [IP]",
	Short:   "get/get node IP address in the meepctl config file",
	Long:    "Get/Set node IP address in the meepctl config file",
	Example: "meepctl config ip 1.2.3.4",
	Args:    cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		key := "node.ip"
		v, _ := cmd.Flags().GetBool("verbose")
		t, _ := cmd.Flags().GetBool("time")
		if v {
			fmt.Println("config ip called")
			fmt.Println("[flag] verbose:", v)
			fmt.Println("[flag] time:", t)
		}

		start := time.Now()
		value := viper.GetString(key)

		if len(args) == 0 {
			_ = cmd.Help()
			fmt.Println("")
		} else {
			ip := args[0]
			valid, reason := utils.ConfigIPValid(ip)
			if valid {
				cfg := utils.ConfigReadFile(viper.ConfigFileUsed())
				cfg.Node.IP = ip
				err := utils.ConfigWriteFile(cfg, viper.ConfigFileUsed())
				if err != nil {
					fmt.Println(err)
				}
				value = ip
			} else {
				fmt.Println("Invalid IP: " + reason)
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
	configCmd.AddCommand(configIp)
}
