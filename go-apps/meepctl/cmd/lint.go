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

// buildCmd represents the build command
var lintCmd = &cobra.Command{
	Use:   "lint <targets>",
	Short: "Lint core components & packages",
	Long: `AdvantEDGE is composed of a collection of micro-services.

Lint command verifies code format & syntax.
Multiple targets can be specified (e.g. meepctl lint <target1> <target2>...)

Valid targets:`,
	Example: `  # Lint all components
  meepctl lint all
  # Lint meep-ctrl-engine component only
  meepctl lint meep-ctrl-engine`,
	Args: cobra.OnlyValidArgs,
	ValidArgs: []string{
		"all",
		"meep-frontend",
		"meep-ctrl-engine",
		"meep-webhook",
		"meep-mg-manager",
		"meep-mon-engine",
		"meep-rnis",
		"meep-loc-serv",
		"meep-metrics-engine",
		"meep-tc-engine",
		"meep-tc-sidecar",
		"meep-virt-engine",
		"meep-couch",
		"meep-ctrl-engine-client",
		"meep-ctrl-engine-model",
		"meep-rnis-client",
		"meep-rnis-notification-client",
		"meep-loc-serv-client",
		"meep-loc-serv-notification-client",
		"meep-logger",
		"meep-http-logger",
		"meep-metric-store",
		"meep-metrics-engine-notification-client",
		"meep-mg-app-client",
		"meep-mg-manager-client",
		"meep-mg-manager-model",
		"meep-model",
		"meep-net-char-mgr",
		"meep-redis",
		"meep-replay-manager",
		"meep-watchdog",
	},

	Run: func(cmd *cobra.Command, args []string) {
		if !utils.ConfigValidate("") {
			fmt.Println("Fix configuration issues")
			return
		}

		targets := args
		if len(targets) == 0 {
			fmt.Println("Need to specify at least one target from ", cmd.ValidArgs)
			os.Exit(0)
		}

		v, _ := cmd.Flags().GetBool("verbose")
		t, _ := cmd.Flags().GetBool("time")

		if v {
			fmt.Println("Lint called")
			fmt.Println("[arg]  targets:", targets)
			fmt.Println("[flag] verbose:", v)
			fmt.Println("[flag] time:", t)
		}

		start := time.Now()
		for _, target := range targets {
			if target == "all" {
				lintAll(cmd)
			} else {
				lint(target, cmd)
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
	for _, arg := range lintCmd.ValidArgs {
		argsStr += "\n  * " + arg
	}
	lintCmd.Long += argsStr
	rootCmd.AddCommand(lintCmd)
}

func lintAll(cobraCmd *cobra.Command) {
	for _, target := range cobraCmd.ValidArgs {
		if target == "all" {
			continue
		}
		lint(target, cobraCmd)
		fmt.Println("")
	}
}

func lint(targetName string, cobraCmd *cobra.Command) {
	target := utils.RepoCfg.GetStringMapString("repo.core." + targetName)
	if len(target) == 0 {
		fmt.Println("Invalid target:", targetName)
		return
	}

	if strings.HasPrefix(target["src"], "go-") {
		lintGo(targetName, cobraCmd)
	} else if strings.HasPrefix(target["src"], "js-") {
		lintJs(targetName, cobraCmd)
	}
}

func lintGo(targetName string, cobraCmd *cobra.Command) {
	fmt.Println("--", targetName, "--")
	target := utils.RepoCfg.GetStringMapString("repo.core." + targetName)
	gitDir := viper.GetString("meep.gitdir")
	srcDir := gitDir + "/" + target["src"]
	lintEnabled := utils.RepoCfg.GetBool("repo.core." + targetName + ".lint")
	if !lintEnabled {
		fmt.Println("   + skipping")
		return
	}

	// linter
	fmt.Println("   + running linter (go)")
	cmd := exec.Command("golangci-lint", "run")
	cmd.Dir = srcDir
	out, err := utils.ExecuteCmd(cmd, cobraCmd)
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println(out)
		fmt.Println("Linting failed. Exiting...")
		return
	}
}

func lintJs(targetName string, cobraCmd *cobra.Command) {
	fmt.Println("--", targetName, "--")
	target := utils.RepoCfg.GetStringMapString("repo.core." + targetName)
	gitDir := viper.GetString("meep.gitdir")
	srcDir := gitDir + "/" + target["src"]
	lintEnabled := utils.RepoCfg.GetBool("repo.core." + targetName + ".lint")
	if !lintEnabled {
		fmt.Println("   + skipping")
		return
	}

	// linter: ESLint
	fmt.Println("   + running linter (js)")
	cmd := exec.Command("eslint", "src/js/")
	cmd.Dir = srcDir
	out, err := utils.ExecuteCmd(cmd, cobraCmd)
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println(out)
		fmt.Println("Linting failed. Exiting...")
		return
	}
}
