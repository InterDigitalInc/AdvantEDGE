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

		targets := utils.RepoCfg.GetStringMapString("repo.core")

		for k := range targets {
			codecovCapable := utils.RepoCfg.GetBool("repo.core." + k + ".codecov")
			if codecovCapable {
				gitDir := viper.GetString("meep.gitdir")
				workDir := viper.GetString("meep.workdir")
				codecovFile := workDir + "/codecov/" + k + "/codecov-" + k + ".out"
				if _, err := os.Stat(codecovFile); !os.IsNotExist(err) {
					fmt.Println("Found " + codecovFile)
					targetDir := gitDir + "/go-apps/" + k

					//go tool cover -html=c.out -o coverage.html
					htmlReport := gitDir + "/test/codecov-" + k + ".html"
					fmt.Println("  + Generating html report ", htmlReport)
					cmd := exec.Command("go", "tool", "cover", "-html="+codecovFile, "-o", htmlReport)
					cmd.Dir = targetDir
					_, _ = utils.ExecuteCmd(cmd, cobraCmd)

					// go tool cover -func=c.out
					txtReport := gitDir + "/test/codecov-" + k + ".txt"
					fmt.Println("  + Generating text report ", txtReport)
					cmd = exec.Command("go", "tool", "cover", "-func="+codecovFile, "-o", txtReport)
					cmd.Dir = targetDir
					_, _ = utils.ExecuteCmd(cmd, cobraCmd)

				}
			}

		}

	},
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
