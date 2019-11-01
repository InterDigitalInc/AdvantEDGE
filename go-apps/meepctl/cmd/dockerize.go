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
	nodeIp := viper.GetString("node.ip")
	apiHost := utils.RepoCfg.GetBool("repo.core." + targetName + ".apihost")

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
				//copy bin/api
				if apiHost {
					dstDataDirApi := dstDataDir + "/api"
					if verbose {
						fmt.Println("    Hosting Api in " + dstDataDirApi)
					}

					//copy swagger-ui files
					cmd = exec.Command("cp", "-r", gitdir+"/swagger-ui/dist", dstDataDir)
					_, err = utils.ExecuteCmd(cmd, cobraCmd)
					if err != nil {
						fmt.Println("Failed to copy: ", err)
						return
					}

					//rename directory
					cmd := exec.Command("mv", dstDataDir+"/dist", dstDataDirApi)
					_, err = utils.ExecuteCmd(cmd, cobraCmd)
					if err != nil {
						fmt.Println("Failed to move: ", err)
						return
					}

					//get all the yaml file to be put in the /api directory as well as putting the host line for the TRY-IT-OUT function to work
					urls := " [ "
					urlStringToReplace := `url: "https:\/\/petstore.swagger.io\/v2\/swagger.json",`

					//find all the apis and copy them at the location above
					for _, targetArg := range cobraCmd.ValidArgs {
						if targetArg == "all" {
							continue
						}
						apiLocationEntry := utils.RepoCfg.GetString("repo.core." + targetArg + ".api")
						if apiLocationEntry == "" {
							continue
						}
						apiLocationFile := gitdir + "/" + apiLocationEntry
						nodePort := utils.RepoCfg.GetString("repo.core." + targetArg + ".nodeport")
						if apiLocationFile != "" {
							dstTargetApiFile := dstDataDirApi + "/" + targetArg + "-api.yaml"
							if verbose {
								fmt.Println("    Copying: " + apiLocationFile + " --> " + dstTargetApiFile)
							}

							cmd = exec.Command("cp", apiLocationFile, dstTargetApiFile)
							_, err = utils.ExecuteCmd(cmd, cobraCmd)
							if err != nil {
								fmt.Println("Failed to copy: ", err)
								return
							}
							//find if host line already exist in the file, if it does, remove it
							cmd = exec.Command("grep", "host: ", dstTargetApiFile)
							hostLine, _ := utils.ExecuteCmd(cmd, cobraCmd)
							if hostLine != "" {
								hostLine = strings.TrimSpace(hostLine)
								sedHostLine := "/" + hostLine + "/d"
								cmd = exec.Command("sed", "-i", sedHostLine, dstTargetApiFile)
								_, err = utils.ExecuteCmd(cmd, cobraCmd)
								if err != nil {
									fmt.Println("Failed to sed: ", err)
									return
								}
							}

							//find the basepath line in the file and append the host line
							cmd = exec.Command("grep", "basePath: ", dstTargetApiFile)
							basePath, err := utils.ExecuteCmd(cmd, cobraCmd)
							if err != nil {
								fmt.Println("Failed to grep: ", err)
								return
							}

							if basePath == "" {
								fmt.Println("Error: basePath shouldn't be empty")
								return
							}

							newHostLine := "host: " + nodeIp + ":" + nodePort
							newBasePath := strings.Replace(basePath, `/`, `\/`, -1)
							//removing the CR/LF at end of line
							newBasePath = newBasePath[:len(newBasePath)-1]

							sedBasePathLine := "/" + newBasePath + "/a" + newHostLine
							cmd = exec.Command("sed", "-i", sedBasePathLine, dstTargetApiFile)
							_, err = utils.ExecuteCmd(cmd, cobraCmd)
							if err != nil {
								fmt.Println("Failed to sed: ", err)
								return
							}

							//update the string to update the drop-down menu in the index.html file of /api
							cmd = exec.Command("grep", "title:", dstTargetApiFile)
							title, err := utils.ExecuteCmd(cmd, cobraCmd)
							if err != nil {
								fmt.Println("Failed to move: ", err)
								return
							}
							title = strings.TrimSpace(title)
							title = title[6:]
							title = strings.TrimSpace(title)

							if title[0] == '"' {
								title = title[1:]
							}
							if title[len(title)-1] == '"' {
								title = title[0 : len(title)-1]
							}

							//update urls for swagger-ui index file
							urls = urls + `{"name": "` + title + `", "url": "` + targetArg + `-api.yaml"},`
						}
					}

					//update swagger-ui index file
					urls = urls + " ],"
					sedString := "s/" + urlStringToReplace + "/urls: " + urls + "/g"
					cmd = exec.Command("sed", "-i", sedString, dstDataDirApi+"/index.html")
					_, err = utils.ExecuteCmd(cmd, cobraCmd)
					if err != nil {
						fmt.Println("Failed to sed: ", err)
						return
					}

				}
			} else {
				fmt.Println("    Source data not found: " + srcDataDir + " --> " + dstDataDir)
			}
		}
	}

	if err != nil {
		fmt.Println("Error dockerizing ", targetName)
	} else {
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
	}
	// cleanup data
	/*	if len(data) != 0 {
			for k := range data {
				dstDataDir := bindir + "/" + k
				cmd := exec.Command("rm", "-r", dstDataDir)
				_, _ = utils.ExecuteCmd(cmd, cobraCmd)
			}
		}
	*/
}
