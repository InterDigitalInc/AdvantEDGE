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
	"os"
	"os/exec"
	"time"

	"github.com/InterDigitalInc/AdvantEDGE/go-apps/meepctl/utils"

	"github.com/roymx/viper"
	"github.com/spf13/cobra"
)

// dockerizeCmd represents the dockerize command
var dockerizeCmd = &cobra.Command{
	Use:   "dockerize <targets>",
	Short: "Dockerize core components",
	Long: `Dockerize core components

AdvantEDGE is composed of a collection of micro-services.

Dockerize command genrates AdvantEDGE Docker images and stores them in
the local Docker registry.
Multiple targets can be specified (e.g. meepctl dockerize <target1> <target2>...)

Valid targets:`,
	Example: `  # Dockerize all components
  meepctl dockerize all
  # Dockerize meep-ctrl-engine component only
  meepctl dockerize meep-ctrl-engine`,
	Args:      cobra.OnlyValidArgs,
	ValidArgs: []string{"all", "meep-ctrl-engine", "meep-webhook", "meep-mg-manager", "meep-mon-engine", "meep-tc-engine", "meep-tc-sidecar"},
	Run: func(cmd *cobra.Command, args []string) {
		if !utils.ConfigValidate("") {
			fmt.Println("Fix configuration issues")
			return
		}

		targets := args
		if len(targets) == 0 {
			fmt.Println("Error: Need to specify at least one target from ", cmd.ValidArgs)
			fmt.Println("")
			_ = cmd.Help()
			return
		}

		r, _ := cmd.Flags().GetString("registry")
		v, _ := cmd.Flags().GetBool("verbose")
		t, _ := cmd.Flags().GetBool("time")

		if v {
			fmt.Println("Dockerize called")
			fmt.Println("[arg]  targets:", targets)
			fmt.Println("[flag] registry:", r)
			fmt.Println("[flag] verbose:", v)
			fmt.Println("[flag] time:", t)
		}

		start := time.Now()
		utils.InitRepoConfig()
		if r == "" {
			r = viper.GetString("meep.registry")
		}
		fmt.Println("Using docker registry:", r)
		for _, target := range targets {
			if target == "all" {
				dockerizeAll(r, cmd)
			} else {
				dockerize(r, target, cmd)
			}
		}

		elapsed := time.Since(start)
		if t {
			fmt.Println("Took ", elapsed.Round(time.Millisecond).String())
		}
	},
}

func init() {
	var argsStr string
	for _, arg := range dockerizeCmd.ValidArgs {
		argsStr += "\n  * " + arg
	}
	dockerizeCmd.Long += argsStr

	rootCmd.AddCommand(dockerizeCmd)
	dockerizeCmd.Flags().StringP("registry", "r", "", "Override registry from config file")
}

func dockerizeAll(registry string, cobraCmd *cobra.Command) {
	for _, target := range cobraCmd.ValidArgs {
		if target == "all" {
			continue
		}
		dockerize(registry, target, cobraCmd)
	}
}

func dockerize(registry string, targetName string, cobraCmd *cobra.Command) {
	verbose, _ := cobraCmd.Flags().GetBool("verbose")
	target := utils.RepoCfg.GetStringMapString("repo.core." + targetName)
	gitdir := viper.GetString("meep.gitdir")
	bindir := gitdir + "/" + target["bin"]

	if len(target) == 0 {
		fmt.Println("Invalid target:", targetName)
		return
	}

	// copy container data locally
	data := utils.RepoCfg.GetStringMapString("repo.core." + targetName + ".docker-data")
	if len(data) != 0 {
		for k, v := range data {
			dstDataDir := bindir + "/" + k
			srcDataDir := gitdir + "/" + v
			if _, err := os.Stat(srcDataDir); !os.IsNotExist(err) {
				if verbose {
					fmt.Println("    Using: " + srcDataDir + " --> " + dstDataDir)
				}
				cmd := exec.Command("rm", "-r", dstDataDir)
				_, _ = utils.ExecuteCmd(cmd, cobraCmd)
				cmd = exec.Command("cp", "-r", srcDataDir, dstDataDir)
				_, _ = utils.ExecuteCmd(cmd, cobraCmd)
			} else {
				fmt.Println("    Source data not found: " + srcDataDir + " --> " + dstDataDir)
			}
		}
	}

	// dockerize & push to private meep docker registry
	path := gitdir + "/" + target["bin"]
	fmt.Println("dockerizing", targetName)
	if registry != "" {
		tag := registry + "/" + targetName
		cmd := exec.Command("docker", "build", "--no-cache", "--rm", "-t", tag, path)
		_, _ = utils.ExecuteCmd(cmd, cobraCmd)
		cmd = exec.Command("docker", "push", tag)
		_, err := utils.ExecuteCmd(cmd, cobraCmd)
		if err != nil {
			fmt.Println("Failed to push", tag, " Error:", err)
			return
		}
	} else {
		cmd := exec.Command("docker", "build", "--no-cache", "--rm", "-t", targetName, path)
		_, _ = utils.ExecuteCmd(cmd, cobraCmd)
	}

	// cleanup data
	if len(data) != 0 {
		for k := range data {
			dstDataDir := bindir + "/" + k
			cmd := exec.Command("rm", "-r", dstDataDir)
			_, _ = utils.ExecuteCmd(cmd, cobraCmd)
		}
	}

}
