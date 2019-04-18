// Copyright © 2019 InterDigital, Inc
// This file is part of {{ .appName }}.

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
	Use:   "dockerize <target>",
	Short: "Dockerize core components",
	Long: `Dockerize core components

AdvantEDGE is composed of a collection of micro-services.

Dockerize command genrates AdvantEDGE Docker images and stores them in
the local Docker registry.

Valid targets:`,
	Example: `  # Dockerize all components
    meepctl dockerize all
  # Dockerize meep-ctrl-engine component only
    meepctl dockerize meep-ctrl-engine
			`,
	Args:      cobra.ExactValidArgs(1),
	ValidArgs: []string{"all", "meep-ctrl-engine", "meep-initializer", "meep-mg-manager", "meep-mon-engine", "meep-tc-engine", "meep-tc-sidecar"},
	Run: func(cmd *cobra.Command, args []string) {
		target := ""
		if len(args) > 0 {
			target = args[0]
		}

		v, _ := cmd.Flags().GetBool("verbose")
		t, _ := cmd.Flags().GetBool("time")

		if v {
			fmt.Println("Dockerize called")
			fmt.Println("[arg]  target:", target)
			fmt.Println("[flag] verbose:", v)
			fmt.Println("[flag] time:", t)
		}

		start := time.Now()
		utils.InitRepoConfig()
		if target == "all" {
			dockerizeAll(cmd)
		} else {
			dockerize(target, cmd)
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dockerizeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dockerizeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func dockerizeAll(cobraCmd *cobra.Command) {
	for _, target := range cobraCmd.ValidArgs {
		if target == "all" {
			continue
		}
		dockerize(target, cobraCmd)
	}
}

func dockerize(targetName string, cobraCmd *cobra.Command) {
	verbose, _ := cobraCmd.Flags().GetBool("verbose")
	target := utils.RepoCfg.GetStringMapString("repo.core." + targetName)
	gitDir := viper.GetString("meep.gitdir")
	binDir := gitDir + "/" + target["bin"]

	if len(target) == 0 {
		fmt.Println("Invalid target:", targetName)
		return
	}

	//copy container data locally
	data := utils.RepoCfg.GetStringMapString("repo.core." + targetName + ".docker-data")
	if len(data) != 0 {
		for k, v := range data {
			dstDataDir := binDir + "/" + k
			srcDataDir := gitDir + "/" + v
			if _, err := os.Stat(srcDataDir); !os.IsNotExist(err) {
				if verbose {
					fmt.Println("    Using: " + srcDataDir + " --> " + dstDataDir)
				}
				cmd := exec.Command("rm", "-r", dstDataDir)
				utils.ExecuteCmd(cmd, cobraCmd)
				cmd = exec.Command("cp", "-r", srcDataDir, dstDataDir)
				utils.ExecuteCmd(cmd, cobraCmd)
			} else {
				fmt.Println("    Source data not found: " + srcDataDir + " --> " + dstDataDir)
			}
		}
	}

	// dockerize
	path := gitDir + "/" + target["bin"]
	fmt.Println("Dockerizing", targetName)
	cmd := exec.Command("docker", "build", "--no-cache", "--rm", "-t", targetName, path)
	utils.ExecuteCmd(cmd, cobraCmd)

	// cleanup data
	if len(data) != 0 {
		for k := range data {
			dstDataDir := binDir + "/" + k
			cmd := exec.Command("rm", "-r", dstDataDir)
			utils.ExecuteCmd(cmd, cobraCmd)
		}
	}

}
