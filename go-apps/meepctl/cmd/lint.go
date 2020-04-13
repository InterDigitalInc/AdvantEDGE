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

type LintData struct {
	gitdir        string
	coreGoApps    []string
	coreJsApps    []string
	sandboxGoApps []string
	goPackages    []string
	jsPackages    []string
}

const lintDesc = `AdvantEDGE is composed of a collection of micro-services.

Lint command verifies code format & syntax.
Multiple targets can be specified (e.g. meepctl lint <target1> <target2>...)`

const lintExample = `  # Lint all components
  meepctl lint all
  # Lint meep-ctrl-engine component only
  meepctl lint meep-ctrl-engine`

// lintCmd represents the lint command
var lintCmd = &cobra.Command{
	Use:       "lint <targets>",
	Short:     "Lint core components & packages",
	Long:      lintDesc,
	Example:   lintExample,
	Args:      cobra.OnlyValidArgs,
	ValidArgs: nil,
	Run:       lintRun,
}

var lintData LintData

func init() {
	// Get targets from repo config file
	lintData.coreGoApps = utils.GetTargets("repo.core.go-apps")
	lintData.coreJsApps = utils.GetTargets("repo.core.js-apps")
	lintData.sandboxGoApps = utils.GetTargets("repo.sandbox.go-apps")
	lintData.goPackages = utils.GetTargets("repo.packages.go-packages")
	lintData.jsPackages = utils.GetTargets("repo.packages.js-packages")

	// Create the list of valid arguments
	baseArgs := []string{"all"}
	configArgs := append(lintData.coreGoApps, lintData.coreJsApps...)
	configArgs = append(configArgs, lintData.sandboxGoApps...)
	configArgs = append(configArgs, lintData.goPackages...)
	configArgs = append(configArgs, lintData.jsPackages...)
	sort.Strings(configArgs)
	lintCmd.ValidArgs = append(baseArgs, configArgs...)

	// Add list of arguments to Example usage
	lintCmd.Example += "\n\nValid Targets:"
	for _, arg := range lintCmd.ValidArgs {
		lintCmd.Example += "\n  * " + arg
	}

	// Add command
	rootCmd.AddCommand(lintCmd)
}

func lintRun(cmd *cobra.Command, args []string) {
	if !utils.ConfigValidate("") {
		fmt.Println("Fix configuration issues")
		return
	}

	targets := args
	if len(targets) == 0 {
		fmt.Println("Error: Need to specify at least one target")
		_ = cmd.Usage()
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

	// Get config
	lintData.gitdir = strings.TrimSuffix(viper.GetString("meep.gitdir"), "/")

	// Lint code
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
}

func lintAll(cobraCmd *cobra.Command) {
	for _, target := range lintData.coreGoApps {
		lintGo(target, "repo.core.go-apps.", cobraCmd)
		fmt.Println("")
	}
	for _, target := range lintData.coreJsApps {
		lintJs(target, "repo.core.js-apps.", cobraCmd)
		fmt.Println("")
	}
	for _, target := range lintData.sandboxGoApps {
		lintGo(target, "repo.sandbox.go-apps.", cobraCmd)
		fmt.Println("")
	}
	for _, target := range lintData.goPackages {
		lintGo(target, "repo.packages.go-packages.", cobraCmd)
		fmt.Println("")
	}
	for _, target := range lintData.jsPackages {
		lintJs(target, "repo.packages.js-packages.", cobraCmd)
		fmt.Println("")
	}
}

func lint(targetName string, cobraCmd *cobra.Command) {
	for _, target := range lintData.coreGoApps {
		if target == targetName {
			lintGo(targetName, "repo.core.go-apps.", cobraCmd)
			return
		}
	}
	for _, target := range lintData.coreJsApps {
		if target == targetName {
			lintJs(targetName, "repo.core.js-apps.", cobraCmd)
			return
		}
	}
	for _, target := range lintData.sandboxGoApps {
		if target == targetName {
			lintGo(targetName, "repo.sandbox.go-apps.", cobraCmd)
			return
		}
	}
	for _, target := range lintData.goPackages {
		if target == targetName {
			lintGo(targetName, "repo.packages.go-packages.", cobraCmd)
			return
		}
	}
	for _, target := range lintData.jsPackages {
		if target == targetName {
			lintJs(targetName, "repo.packages.js-packages.", cobraCmd)
			return
		}
	}
	fmt.Println("Error: Unsupported target: ", targetName)
}

func lintGo(targetName string, repo string, cobraCmd *cobra.Command) {
	fmt.Println("--", targetName, "--")
	target := utils.RepoCfg.GetStringMapString(repo + targetName)
	srcDir := lintData.gitdir + "/" + target["src"]
	lintEnabled := utils.RepoCfg.GetBool(repo + targetName + ".lint")
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

func lintJs(targetName string, repo string, cobraCmd *cobra.Command) {
	fmt.Println("--", targetName, "--")
	target := utils.RepoCfg.GetStringMapString(repo + targetName)
	srcDir := lintData.gitdir + "/" + target["src"]
	lintEnabled := utils.RepoCfg.GetBool(repo + targetName + ".lint")
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
