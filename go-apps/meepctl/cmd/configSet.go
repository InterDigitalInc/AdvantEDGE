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

var ip, gitdir string // flags

// configSet represents the set command
var configSet = &cobra.Command{
	Use:   "set",
	Short: "Set value(s) in the meepctl config file",
	Long: `Set value(s) in the meepctl config file
	`,
	Example: `  # Configure IP address
    meepctl config set --ip 1.2.3.4
  # Configure GIT directory
    meepctl config set --gitdir /home/some-user/AdvantEDGE
  # Configure GIT to local directory + IP simultaneously
    meepctl config set --gitdir /home/some-user/AdvantEDGE --ip 1.2.3.4
	`,
	Args:      cobra.OnlyValidArgs,
	ValidArgs: []string{"ip", "gitdir"},
	Run: func(cmd *cobra.Command, args []string) {
		v, _ := cmd.Flags().GetBool("verbose")
		t, _ := cmd.Flags().GetBool("time")
		if v {
			fmt.Println("config set called")
			fmt.Println("[flag] ip:", ip)
			fmt.Println("[flag] gitdir:", gitdir)
			fmt.Println("[flag] verbose:", v)
			fmt.Println("[flag] time:", t)
		}

		fmt.Println("len(args)", len(args), "|| args ", args)

		start := time.Now()
		if ip != "" {
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
			}
		}

		if gitdir != "" {
			valid, reason := utils.ConfigGitdirValid(gitdir)
			if valid {
				cfg := utils.ConfigReadFile(viper.ConfigFileUsed())
				cfg.Meep.Gitdir = gitdir
				err := utils.ConfigWriteFile(cfg, viper.ConfigFileUsed())
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("Updated meep.gitdir with [" + gitdir + "]")
			} else {
				fmt.Println("Invalid Gitdir: " + reason)
			}
		}

		if ip == "" && gitdir == "" {
			fmt.Println("Which flag do you want to set", cmd.ValidArgs, "?")
			fmt.Println()
			fmt.Println(cmd.Long)
			fmt.Println(cmd.Example)
		}

		elapsed := time.Since(start)
		if t {
			fmt.Println("Took ", elapsed.Round(time.Millisecond).String())
		}

	},
}

func init() {
	configCmd.AddCommand(configSet)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	configSet.Flags().StringVar(&ip, "ip", "", "IP address of AdvantEDGE node (local IP address)")
	configSet.Flags().StringVar(&gitdir, "gitdir", "", "Path to the AdvantEDGE GIT directory")

}
