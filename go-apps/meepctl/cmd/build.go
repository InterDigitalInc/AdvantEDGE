/*
 * Copyright (c) 2022  The AdvantEDGE Authors
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
	"time"

	"github.com/InterDigitalInc/AdvantEDGE/go-apps/meepctl/utils"

	"github.com/roymx/viper"
	"github.com/spf13/cobra"
)

type BuildData struct {
	codecov       bool
	nolint        bool
	coreGoApps    []string
	coreJsApps    []string
	sandboxGoApps []string
}

const buildDesc = `AdvantEDGE is composed of a collection of micro-services.

Build command generates AdvantEDGE binaries.
Multiple targets can be specified (e.g. meepctl build <target1> <target2>...)`

const buildExample = `  # Build all components
  meepctl build all
  # Build meep-platform-ctrl component only
  meepctl build meep-platform-ctrl`

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:       "build <targets>",
	Short:     "Build core components",
	Long:      buildDesc,
	Example:   buildExample,
	Args:      cobra.OnlyValidArgs,
	ValidArgs: nil,
	Run:       buildRun,
}

var buildData BuildData

func init() {
	// Get targets from repo config file
	buildData.coreGoApps = utils.GetTargets("repo.core.go-apps", "build")
	buildData.coreJsApps = utils.GetTargets("repo.core.js-apps", "build")
	buildData.sandboxGoApps = utils.GetTargets("repo.sandbox.go-apps", "build")

	// Create the list of valid arguments
	baseArgs := []string{"all"}
	configArgs := append(buildData.coreGoApps, buildData.coreJsApps...)
	configArgs = append(configArgs, buildData.sandboxGoApps...)
	sort.Strings(configArgs)
	buildCmd.ValidArgs = append(baseArgs, configArgs...)

	// Add list of arguments to Example usage
	buildCmd.Example += "\n\nValid Targets:"
	for _, arg := range buildCmd.ValidArgs {
		buildCmd.Example += "\n  * " + arg
	}

	// Set build-specific flags
	buildCmd.Flags().BoolVar(&buildData.codecov, "codecov", false, "Build a code coverage binary (dev. option)")
	buildCmd.Flags().BoolVar(&buildData.nolint, "nolint", false, "Disable linting")

	// Add command
	rootCmd.AddCommand(buildCmd)
}

func buildRun(cmd *cobra.Command, args []string) {
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
		fmt.Println("Build called")
		fmt.Println("[arg]  targets:", targets)
		fmt.Println("[flag] codecov:", buildData.codecov)
		fmt.Println("[flag] nolint:", buildData.nolint)
		fmt.Println("[flag] verbose:", v)
		fmt.Println("[flag] time:", t)
	}

	start := time.Now()
	for _, target := range targets {
		if target == "all" {
			buildAll(cmd)
		} else {
			build(target, cmd)
		}
	}
	elapsed := time.Since(start)
	if t {
		fmt.Println("Took ", elapsed.Round(time.Millisecond).String())
	}
}

func buildAll(cobraCmd *cobra.Command) {
	for _, target := range buildData.coreGoApps {
		buildGoApp(target, "repo.core.go-apps.", cobraCmd)
		fmt.Println("")
	}
	for _, target := range buildData.coreJsApps {
		buildJsApp(target, "repo.core.js-apps.", cobraCmd)
		fmt.Println("")
	}
	for _, target := range buildData.sandboxGoApps {
		buildGoApp(target, "repo.sandbox.go-apps.", cobraCmd)
		fmt.Println("")
	}
}

func build(targetName string, cobraCmd *cobra.Command) {
	for _, target := range buildData.coreGoApps {
		if target == targetName {
			buildGoApp(target, "repo.core.go-apps.", cobraCmd)
			return
		}
	}
	for _, target := range buildData.coreJsApps {
		if target == targetName {
			buildJsApp(target, "repo.core.js-apps.", cobraCmd)
			return
		}
	}
	for _, target := range buildData.sandboxGoApps {
		if target == targetName {
			buildGoApp(target, "repo.sandbox.go-apps.", cobraCmd)
			return
		}
	}
	fmt.Println("Error: Unsupported target: ", targetName)
}

func buildJsApp(targetName string, repo string, cobraCmd *cobra.Command) {
	switch targetName {
	case "meep-frontend":
		buildFrontend(targetName, repo, cobraCmd)
	default:
		fmt.Println("Error: Unsupported JS App: ", targetName)
	}
}

func buildFrontend(targetName string, repo string, cobraCmd *cobra.Command) {
	fmt.Println("--", targetName, "--")
	target := utils.RepoCfg.GetStringMapString(repo + targetName)
	version := utils.RepoCfg.GetString("version")
	gitDir := viper.GetString("meep.gitdir")
	srcDir := gitDir + "/" + target["src"]
	binDir := gitDir + "/" + target["bin"]
	lintEnabled := utils.RepoCfg.GetBool(repo + targetName + ".lint")

	// dependencies
	fmt.Println("   + checking external dependencies")
	cmd := exec.Command("npm", "ci")
	cmd.Dir = srcDir
	out, err := utils.ExecuteCmd(cmd, cobraCmd)
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println(out)
	}

	locDeps := utils.RepoCfg.GetStringMapString(repo + targetName + ".local-deps")
	if len(locDeps) > 0 {
		fmt.Println("   + checking local dependencies")
		for dep, depDir := range locDeps {
			fmt.Println("     * " + dep)
			cmd := exec.Command("npm", "ci")
			cmd.Dir = gitDir + "/" + depDir
			out, err := utils.ExecuteCmd(cmd, cobraCmd)
			if err != nil {
				fmt.Println("Error:", err)
				fmt.Println(out)
			}
		}
	}

	// remove old binDir if exists
	if _, err := os.Stat(binDir); !os.IsNotExist(err) {
		cmd = exec.Command("rm", "-r", binDir)
		cmd.Dir = srcDir
		out, err = utils.ExecuteCmd(cmd, cobraCmd)
		if err != nil {
			fmt.Println("Error:", err)
			fmt.Println(out)
		}
	}

	// linter: ESLint
	if lintEnabled && !buildData.nolint {
		fmt.Println("   + running linter")
		cmd := exec.Command("eslint", "src/js/")
		cmd.Dir = srcDir
		out, err := utils.ExecuteCmd(cmd, cobraCmd)
		if err != nil {
			fmt.Println("Error:", err)
			fmt.Println(out)
			fmt.Println("Linting failed. Exiting...")
			fmt.Println("To skip linting run build with --nolint")
			return
		}
	}

	//build
	fmt.Println("   + building " + targetName)
	cmd = exec.Command("npm", "run", "build", "--", "--output-path="+binDir, "--env.VERSION=v"+version)
	cmd.Dir = srcDir
	out, err = utils.ExecuteCmd(cmd, cobraCmd)
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println(out)
	}
}

func buildGoApp(targetName string, repo string, cobraCmd *cobra.Command) {
	fmt.Println("--", targetName, "--")
	target := utils.RepoCfg.GetStringMapString(repo + targetName)
	gitDir := viper.GetString("meep.gitdir")
	srcDir := gitDir + "/" + target["src"]
	binDir := gitDir + "/" + target["bin"]
	codecovCapable := utils.RepoCfg.GetBool(repo + targetName + ".codecov")
	lintEnabled := utils.RepoCfg.GetBool(repo + targetName + ".lint")

	// dependencies
	fmt.Println("   + checking external dependencies")
	cmd := exec.Command("go", "mod", "vendor")
	cmd.Dir = srcDir
	out, err := utils.ExecuteCmd(cmd, cobraCmd)
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println(out)
	}

	// linter: goloangci-lint
	if lintEnabled && !buildData.nolint {
		fmt.Println("   + running linter")
		cmd := exec.Command("golangci-lint", "run")
		cmd.Dir = srcDir
		out, err := utils.ExecuteCmd(cmd, cobraCmd)
		if err != nil {
			fmt.Println("Error:", err)
			fmt.Println(out)
			fmt.Println("Linting failed. Exiting...")
			fmt.Println("To skip linting run build with --nolint")
			return
		}
	}

	// for go, local dependencies are handled via the Go toolchain so nothing to do

	// Remove unnecessary deps
	fixDeps(targetName, repo, cobraCmd)

	// build
	buildFlags := utils.RepoCfg.GetStringSlice(repo + targetName + ".build-flags")
	if buildData.codecov && codecovCapable {
		fmt.Println("   + building " + targetName + " (warning: development build - code coverage)")
		args := []string{"test", "-covermode=count", "-coverpkg=./...", "-c"}
		args = append(args, buildFlags...)
		args = append(args, "-o", binDir+"/"+targetName, srcDir)
		cmd = exec.Command("go", args...)
	} else {
		fmt.Println("   + building " + targetName)
		args := []string{"build"}
		args = append(args, buildFlags...)
		args = append(args, "-o", binDir+"/"+targetName, srcDir)
		cmd = exec.Command("go", args...)
	}
	cmd.Dir = srcDir
	out, err = utils.ExecuteCmd(cmd, cobraCmd)
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println(out)
	}
}

func fixDeps(targetName string, repo string, cobraCmd *cobra.Command) {
	target := utils.RepoCfg.GetStringMapString(repo + targetName)
	gitDir := viper.GetString("meep.gitdir")
	srcDir := gitDir + "/" + target["src"]

	switch targetName {
	case "meep-webhook":
		cmd := exec.Command("rm", "-Rf", srcDir+"/vendor/github.com/hashicorp/golang-lru")
		out, err := utils.ExecuteCmd(cmd, cobraCmd)
		if err != nil {
			fmt.Println("Error:", err)
			fmt.Println(out)
		}

	case "meep-mon-engine":
		cmd := exec.Command("rm", srcDir+"/vendor/k8s.io/client-go/tools/cache/mutation_cache.go")
		out, err := utils.ExecuteCmd(cmd, cobraCmd)
		if err != nil {
			fmt.Println("Error:", err)
			fmt.Println(out)
		}
		cmd = exec.Command("rm", "-Rf", srcDir+"/vendor/github.com/hashicorp/golang-lru")
		out, err = utils.ExecuteCmd(cmd, cobraCmd)
		if err != nil {
			fmt.Println("Error:", err)
			fmt.Println(out)
		}
	}
}
