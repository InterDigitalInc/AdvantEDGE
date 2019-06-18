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
var configGitdir = &cobra.Command{
	Use:     "gitdir [GIT dir path]",
	Short:   "get/set GIT directory path in the meepctl config file",
	Long:    "Get/Set GIT directory path in the meepctl config file",
	Example: "  meepctl config gitdir /home/some-user/AdvantEDGE",
	Args:    cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		key := "meep.gitdir"
		v, _ := cmd.Flags().GetBool("verbose")
		t, _ := cmd.Flags().GetBool("time")
		if v {
			fmt.Println("config gitdir called")
			fmt.Println("[flag] verbose:", v)
			fmt.Println("[flag] time:", t)
		}

		start := time.Now()
		value := viper.GetString(key)
		if len(args) == 0 {
			_ = cmd.Help()
			fmt.Println("")
		} else {
			gitdir := args[0]
			valid, reason := utils.ConfigPathValid(gitdir)
			if valid {
				cfg := utils.ConfigReadFile(viper.ConfigFileUsed())
				cfg.Meep.Gitdir = gitdir
				err := utils.ConfigWriteFile(cfg, viper.ConfigFileUsed())
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("Updated meep.gitdir with [" + gitdir + "]")
				value = gitdir
			} else {
				fmt.Println("Invalid Gitdir: " + reason)
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
	configCmd.AddCommand(configGitdir)
}
