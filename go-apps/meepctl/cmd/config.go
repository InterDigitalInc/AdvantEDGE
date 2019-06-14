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
		fmt.Println("CONFIGURED VALUES")
		for _, key := range keys {
			fmt.Println(key, ":", viper.GetString(key))
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
