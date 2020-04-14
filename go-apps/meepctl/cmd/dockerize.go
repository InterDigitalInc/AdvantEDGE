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
	"sort"
	"strings"
	"time"

	"github.com/InterDigitalInc/AdvantEDGE/go-apps/meepctl/utils"

	"github.com/roymx/viper"
	"github.com/spf13/cobra"
)

type DockerizeData struct {
	registry      string
	gitdir        string
	coreGoApps    []string
	sandboxGoApps []string
}

const dockerizeDesc = `Dockerize core components

AdvantEDGE is composed of a collection of micro-services.

Dockerize command genrates AdvantEDGE Docker images and stores them in
the local Docker registry.
Multiple targets can be specified (e.g. meepctl dockerize <target1> <target2>...)`

const dockerizeExample = `  # Dockerize all components
  meepctl dockerize all
  # Dockerize meep-ctrl-engine component only
  meepctl dockerize meep-ctrl-engine`

// dockerizeCmd represents the dockerize command
var dockerizeCmd = &cobra.Command{
	Use:       "dockerize <targets>",
	Short:     "Dockerize core components",
	Long:      dockerizeDesc,
	Example:   dockerizeExample,
	Args:      cobra.OnlyValidArgs,
	ValidArgs: nil,
	Run:       dockerizeRun,
}

var dockerizeData DockerizeData

func init() {
	// Get targets from repo config file
	dockerizeData.coreGoApps = utils.GetTargets("repo.core.go-apps")
	dockerizeData.sandboxGoApps = utils.GetTargets("repo.sandbox.go-apps")

	// Create the list of valid arguments
	baseArgs := []string{"all"}
	configArgs := append(dockerizeData.coreGoApps, dockerizeData.sandboxGoApps...)
	sort.Strings(configArgs)
	dockerizeCmd.ValidArgs = append(baseArgs, configArgs...)

	// Add list of arguments to Example usage
	dockerizeCmd.Example += "\n\nValid Targets:"
	for _, arg := range dockerizeCmd.ValidArgs {
		dockerizeCmd.Example += "\n  * " + arg
	}

	// Set dockerize-specific flags
	dockerizeCmd.Flags().StringP("registry", "r", "", "Override registry from config file")

	// Add command
	rootCmd.AddCommand(dockerizeCmd)
}

func dockerizeRun(cmd *cobra.Command, args []string) {
	if !utils.ConfigValidate("") {
		fmt.Println("Fix configuration issues")
		return
	}

	targets := args
	if len(targets) == 0 {
		fmt.Println("Error: Need to specify at least one target")
		_ = cmd.Usage()
		return
	}

	dockerizeData.registry, _ = cmd.Flags().GetString("registry")
	v, _ := cmd.Flags().GetBool("verbose")
	t, _ := cmd.Flags().GetBool("time")

	if v {
		fmt.Println("Dockerize called")
		fmt.Println("[arg]  targets:", targets)
		fmt.Println("[flag] registry:", dockerizeData.registry)
		fmt.Println("[flag] verbose:", v)
		fmt.Println("[flag] time:", t)
	}

	start := time.Now()

	// Retrieve registry from config file if not already set
	if dockerizeData.registry == "" {
		dockerizeData.registry = viper.GetString("meep.registry")
	}
	dockerizeData.registry = strings.TrimSuffix(dockerizeData.registry, "/")
	fmt.Println("Using docker registry:", dockerizeData.registry)

	// Get config
	dockerizeData.gitdir = strings.TrimSuffix(viper.GetString("meep.gitdir"), "/")

	// Dockerize microservices
	for _, target := range targets {
		if target == "all" {
			dockerizeAll(cmd)
		} else {
			dockerizeTarget(target, cmd)
		}
	}

	elapsed := time.Since(start)
	if t {
		fmt.Println("Took ", elapsed.Round(time.Millisecond).String())
	}
}

func dockerizeAll(cobraCmd *cobra.Command) {
	for _, target := range dockerizeData.coreGoApps {
		dockerize(target, "repo.core.go-apps.", cobraCmd)
		fmt.Println("")
	}
	for _, target := range dockerizeData.sandboxGoApps {
		dockerize(target, "repo.sandbox.go-apps.", cobraCmd)
		fmt.Println("")
	}
}

func dockerizeTarget(targetName string, cobraCmd *cobra.Command) {
	for _, target := range dockerizeData.coreGoApps {
		if target == targetName {
			dockerize(target, "repo.core.go-apps.", cobraCmd)
			return
		}
	}
	for _, target := range dockerizeData.sandboxGoApps {
		if target == targetName {
			dockerize(target, "repo.sandbox.go-apps.", cobraCmd)
			return
		}
	}
	fmt.Println("Error: Unsupported target: ", targetName)
}

func dockerize(targetName string, repo string, cobraCmd *cobra.Command) {
	verbose, _ := cobraCmd.Flags().GetBool("verbose")
	bindir := dockerizeData.gitdir + "/" + utils.RepoCfg.GetString(repo+targetName+".bin")
	fmt.Println("--", targetName, "--")

	// copy container data locally
	fmt.Println("   + copy docker data")
	data := utils.RepoCfg.GetStringMapString(repo + targetName + ".docker-data")
	var err error
	if len(data) != 0 {
		for k, v := range data {
			dstDataDir := bindir + "/" + k
			srcDataDir := dockerizeData.gitdir + "/" + v
			if _, err = os.Stat(srcDataDir); !os.IsNotExist(err) {
				if verbose {
					fmt.Println("    Using: " + srcDataDir + " --> " + dstDataDir)
				}
				cmd := exec.Command("rm", "-r", dstDataDir)
				_, _ = utils.ExecuteCmd(cmd, cobraCmd)
				cmd = exec.Command("cp", "-r", srcDataDir, dstDataDir)
				_, err = utils.ExecuteCmd(cmd, cobraCmd)
				if err != nil {
					fmt.Println("Error: Failed to copy data: ", srcDataDir, " --> ", dstDataDir)
					return
				}
			} else {
				fmt.Println("Error: Source data not found: ", srcDataDir, " --> ", dstDataDir)
				return
			}
		}
	}

	// Obtain checksum of bin folder contents to add as a label in docker image
	cmd := exec.Command("/bin/sh", "-c", "find "+bindir+" -type f | xargs sha256sum | sort | sha256sum")
	output, _ := utils.ExecuteCmd(cmd, cobraCmd)
	checksum := strings.Split(output, " ")

	// dockerize & push to private meep docker registry
	fmt.Println("   + dockerize " + targetName)
	if dockerizeData.registry != "" {
		tag := dockerizeData.registry + "/" + targetName
		cmd := exec.Command("docker", "build", "--no-cache", "--rm", "--label", "MeepVersion="+checksum[0], "-t", tag, bindir)
		_, err = utils.ExecuteCmd(cmd, cobraCmd)
		if err != nil {
			fmt.Println("Error: Failed to dockerize ", tag, " with error: ", err)
			return
		}
		cmd = exec.Command("docker", "push", tag)
		_, err := utils.ExecuteCmd(cmd, cobraCmd)
		if err != nil {
			fmt.Println("Error: Failed to push ", tag, " with error: ", err)
			return
		}
	} else {
		cmd := exec.Command("docker", "build", "--no-cache", "--rm", "--label", "MeepVersion="+checksum[0], "-t", targetName, bindir)
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
