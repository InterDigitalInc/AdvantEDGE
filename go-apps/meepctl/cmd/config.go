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

	"github.com/roymx/viper"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Config allows to manage meepctl configuration",
	Long: `Config allows to manage meepctl configuration

meepctl relies on a configuration file that lives here ~/.meepctl.yaml

On first meepctl execution, the configuration file is created with default values
It then needs to be initialized once by running initial configuration command (below)
`,
	Example: ` # Initial configuration
 meepctl config set --ip <your-node-ip> --gitdir <path-to-advantedge-git-dir>
 # Help on set command
 meepctl config set --help`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("")
		fmt.Println("   PARAMETER\t\tDESCRIPTION\t\t\t CURRENT VALUE")
		fmt.Println("   version\t\tconfig file version\t\t", viper.GetString("version"))
		fmt.Println("   node.ip\t\tnode's IP address\t\t", viper.GetString("node.ip"))
		fmt.Println("   meep.gitdir\t\tAdvantEDGE repo path\t\t", viper.GetString("meep.gitdir"))
		fmt.Println("   meep.workdir\t\tRuntime storage path\t\t", viper.GetString("meep.workdir"))
		fmt.Println("")
		fmt.Println(cmd.Long)
		fmt.Println(cmd.Example)
		fmt.Println("")
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
