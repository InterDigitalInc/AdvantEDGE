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

	"github.com/InterDigitalInc/AdvantEDGE/go-apps/meepctl/utils"

	"github.com/roymx/viper"
	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Generate code coverage report",
	Long: `Generate code coverage report

	AdvantEDGE components can be compiled in code coverage (codecov) mode.
	When instrumented with codecov, AdventEDGE micro-services generate a coverage report on termination.

	Use this command after terminating codecov execution to genrate a report.
	`,
	Run: func(cobraCmd *cobra.Command, args []string) {
		if !utils.ConfigValidate("") {
			fmt.Println("Fix configuration issues")
			return
		}

		platformTargets := utils.RepoCfg.GetStringMapString("repo.core.go-apps")
		sandboxTargets := utils.RepoCfg.GetStringMapString("repo.sandbox.go-apps")

		for k := range platformTargets {
			codecovCapable := utils.RepoCfg.GetBool("repo.core.go-apps." + k + ".codecov")
			if codecovCapable {
				getCoverageReports(k, cobraCmd, args)
			}
		}
		for k := range sandboxTargets {
			codecovCapable := utils.RepoCfg.GetBool("repo.sandbox.go-apps." + k + ".codecov")
			if codecovCapable {
				getCoverageReports(k, cobraCmd, args)

			}
		}
	},
}

func getCoverageReports(k string, cobraCmd *cobra.Command, args []string) {
	gitDir := viper.GetString("meep.gitdir")
	workDir := viper.GetString("meep.workdir") + "/codecov/" + k
	targetDir := gitDir + "/go-apps/" + k
	codecovFile := gitDir + "/test/codecov/" + k + "-aggregated.out"
	f, err := os.Open(workDir)
	if err != nil {
		fmt.Println(err)
		return
	}
	files, err := f.Readdir(0)
	if err != nil {
		fmt.Println(err)
		return
	}
	cmdArgs := []string{"-cc", codecovFile}
	for _, v := range files {
		name := workDir + "/" + v.Name()
		cmdArgs = append(cmdArgs, name)
	}
	cmd := exec.Command("cov-report", cmdArgs...)
	cmd.Dir = targetDir
	_, err = utils.ExecuteCmd(cmd, cobraCmd)
	if err != nil {
		fmt.Println(err.Error())
	}
	if _, err := os.Stat(codecovFile); !os.IsNotExist(err) {
		fmt.Println("Found " + codecovFile)

		//go tool cover -html=c.out -o coverage.html
		htmlReport := gitDir + "/test/codecov/" + k + "-aggregated.html"
		fmt.Println("  + Generating html report ", htmlReport)
		cmd := exec.Command("go", "tool", "cover", "-html="+codecovFile, "-o", htmlReport)
		cmd.Dir = targetDir
		_, _ = utils.ExecuteCmd(cmd, cobraCmd)

		// go tool cover -func=c.out
		txtReport := gitDir + "/test/codecov/" + k + "-aggregated.txt"
		fmt.Println("  + Generating text report ", txtReport)
		cmd = exec.Command("go", "tool", "cover", "-func="+codecovFile, "-o", txtReport)
		cmd.Dir = targetDir
		_, _ = utils.ExecuteCmd(cmd, cobraCmd)
	} else {
		fmt.Println(err)
	}
}

func init() {
	rootCmd.AddCommand(testCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
