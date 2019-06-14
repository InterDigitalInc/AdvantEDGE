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
var configRegistry = &cobra.Command{
	Use:   "registry [name:port]",
	Short: "get/get docker registry meepctl config file",
	Long:  `Get/Set docker registry in the meepctl config file`,
	Example: `  # Get currently configured docker registry
  meepctl config registry
  # Set docker registry
  meepctl config registry meep-docker-registry:30001`,
	Run: func(cmd *cobra.Command, args []string) {
		v, _ := cmd.Flags().GetBool("verbose")
		t, _ := cmd.Flags().GetBool("time")
		if v {
			fmt.Println("config registry called")
			fmt.Println("[flag] verbose:", v)
			fmt.Println("[flag] time:", t)
		}

		start := time.Now()
		if len(args) > 0 {
			registry := args[0]
			// valid, reason := utils.ConfigGitdirValid(gitdir)
			// if valid {
			cfg := utils.ConfigReadFile(viper.ConfigFileUsed())
			cfg.Meep.Registry = registry
			err := utils.ConfigWriteFile(cfg, viper.ConfigFileUsed())
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Updated meep.registry with [" + registry + "]")
			// } else {
			// 	fmt.Println("Invalid Gitdir: " + reason)
			// 	fmt.Println("")
			// 	_ = cmd.Help()
			// }
		} else {
			key := "meep.registry"
			fmt.Println(key, ":", viper.GetString(key))
		}

		elapsed := time.Since(start)
		if t {
			fmt.Println("Took ", elapsed.Round(time.Millisecond).String())
		}
	},
}

func init() {
	configCmd.AddCommand(configRegistry)
}
