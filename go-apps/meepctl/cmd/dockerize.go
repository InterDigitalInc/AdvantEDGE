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
	"os"
	"os/exec"
	"strings"
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
	ValidArgs: []string{"all", "meep-ctrl-engine", "meep-webhook", "meep-mg-manager", "meep-mon-engine", "meep-loc-serv", "meep-metrics-engine", "meep-tc-engine", "meep-tc-sidecar"},
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
	var err error
	if len(data) != 0 {
		for k, v := range data {
			dstDataDir := bindir + "/" + k
			srcDataDir := gitdir + "/" + v
			if _, err = os.Stat(srcDataDir); !os.IsNotExist(err) {
				if verbose {
					fmt.Println("    Using: " + srcDataDir + " --> " + dstDataDir)
				}
				cmd := exec.Command("rm", "-r", dstDataDir)
				_, _ = utils.ExecuteCmd(cmd, cobraCmd)
				cmd = exec.Command("cp", "-r", srcDataDir, dstDataDir)
				_, err = utils.ExecuteCmd(cmd, cobraCmd)
			} else {
				fmt.Println("    Source data not found: " + srcDataDir + " --> " + dstDataDir)
			}
		}
	}

	if err == nil {
		// Obtain checksum of bin folder contents to add as a label in docker image
		path := gitdir + "/" + target["bin"]
		cmd := exec.Command("/bin/sh", "-c", "find "+path+" -type f | xargs sha256sum | sort | sha256sum")
		output, _ := utils.ExecuteCmd(cmd, cobraCmd)
		checksum := strings.Split(output, " ")

		// dockerize & push to private meep docker registry
		fmt.Println("dockerizing", targetName)
		if registry != "" {
			tag := registry + "/" + targetName
			cmd := exec.Command("docker", "build", "--no-cache", "--rm", "--label", "MeepVersion="+checksum[0], "-t", tag, path)
			_, _ = utils.ExecuteCmd(cmd, cobraCmd)
			cmd = exec.Command("docker", "push", tag)
			_, err := utils.ExecuteCmd(cmd, cobraCmd)
			if err != nil {
				fmt.Println("Failed to push", tag, " Error:", err)
				return
			}
		} else {
			cmd := exec.Command("docker", "build", "--no-cache", "--rm", "--label", "MeepVersion="+checksum[0], "-t", targetName, path)
			_, _ = utils.ExecuteCmd(cmd, cobraCmd)
		}
	} else {
		fmt.Println("dockerizing could not be initiated: Error with the build for", targetName)
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
