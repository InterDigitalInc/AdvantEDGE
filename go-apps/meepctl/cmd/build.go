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

var buildCodecov bool
var buildNolint bool

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build <targets>",
	Short: "Build core components",
	Long: `AdvantEDGE is composed of a collection of micro-services.

Build command generates AdvantEDGE binaries.
Multiple targets can be specified (e.g. meepctl build <target1> <target2>...)

Valid targets:`,
	Example: `  # Build all components
  meepctl build all
  # Build meep-ctrl-engine component only
  meepctl build meep-ctrl-engine`,
	Args: cobra.OnlyValidArgs,
	// WARNING -- meep-frontend comes before meep-ctrl-engine so that "all" works
	ValidArgs: []string{"all", "meep-frontend", "meep-ctrl-engine", "meep-swagger-ui", "meep-webhook", "meep-mg-manager", "meep-mon-engine", "meep-loc-serv", "meep-metrics-engine", "meep-tc-engine", "meep-tc-sidecar", "meep-virt-engine"},

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
			fmt.Println("Build called")
			fmt.Println("[arg]  targets:", targets)
			fmt.Println("[flag] codecov:", buildCodecov)
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
	},
}

func init() {
	var argsStr string
	for _, arg := range buildCmd.ValidArgs {
		argsStr += "\n  * " + arg
	}
	buildCmd.Long += argsStr

	rootCmd.AddCommand(buildCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	buildCmd.Flags().BoolVar(&buildCodecov, "codecov", false, "Build a code coverage binary (dev. option)")
	buildCmd.Flags().BoolVar(&buildNolint, "nolint", false, "Disable linting")
}

func buildAll(cobraCmd *cobra.Command) {
	for _, target := range cobraCmd.ValidArgs {
		if target == "all" {
			continue
		}
		build(target, cobraCmd)
		fmt.Println("")
	}
}

func build(targetName string, cobraCmd *cobra.Command) {
	target := utils.RepoCfg.GetStringMapString("repo.core." + targetName)
	if len(target) == 0 {
		fmt.Println("Invalid target:", targetName)
		return
	}

	switch targetName {
	case "meep-frontend":
		buildFrontend(targetName, cobraCmd)
	case "meep-swagger-ui":
		buildSwaggerUi(targetName, cobraCmd)
	default:
		buildGoApp(targetName, cobraCmd)
	}
}

func buildFrontend(targetName string, cobraCmd *cobra.Command) {
	fmt.Println("--", targetName, "--")
	target := utils.RepoCfg.GetStringMapString("repo.core." + targetName)
	gitDir := viper.GetString("meep.gitdir")
	srcDir := gitDir + "/" + target["src"]
	binDir := gitDir + "/" + target["bin"]
	lintEnabled := utils.RepoCfg.GetBool("repo.core." + targetName + ".lint")

	// dependencies
	fmt.Println("   + checking external dependencies")
	cmd := exec.Command("npm", "ci")
	cmd.Dir = srcDir
	out, err := utils.ExecuteCmd(cmd, cobraCmd)
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println(out)
	}

	locDeps := utils.RepoCfg.GetStringMapString("repo.core." + targetName + ".local-deps")
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
	if lintEnabled && !buildNolint {
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
	cmd = exec.Command("npm", "run", "build", "--", "--output-path="+binDir)
	cmd.Dir = srcDir
	out, err = utils.ExecuteCmd(cmd, cobraCmd)
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println(out)
	}
}

func buildSwaggerUi(targetName string, cobraCmd *cobra.Command) {
	fmt.Println("--", targetName, "--")
	fmt.Println("   + creating api files")

	verbose, _ := cobraCmd.Flags().GetBool("verbose")
	target := utils.RepoCfg.GetStringMapString("repo.core." + targetName)
	gitDir := viper.GetString("meep.gitdir")
	srcDir := gitDir + "/" + target["src"]
	binDir := gitDir + "/" + target["bin"]
	nodeIp := viper.GetString("node.ip")
	ctrlEnginePort := utils.RepoCfg.GetString("repo.core.meep-ctrl-engine.nodeport")

	// remove old binDir if exists
	if _, err := os.Stat(binDir); !os.IsNotExist(err) {
		cmd := exec.Command("rm", "-r", binDir)
		cmd.Dir = srcDir
		out, err := utils.ExecuteCmd(cmd, cobraCmd)
		if err != nil {
			fmt.Println("Error:", err)
			fmt.Println(out)
		}
	}

	if verbose {
		fmt.Println("    Copy hosting Api in " + binDir + " from " + srcDir)
	}
	//copy swagger-ui files
	cmd := exec.Command("cp", "-r", srcDir, binDir)
	_, err := utils.ExecuteCmd(cmd, cobraCmd)
	if err != nil {
		fmt.Println("Failed to copy: ", err)
		return
	}

	//get all the yaml file to be put in the /api directory as well as putting the host line for the TRY-IT-OUT function to work
	urls := " [ "
	urlStringToReplace := `url: "https:\/\/petstore.swagger.io\/v2\/swagger.json",`

	apiBundle := utils.RepoCfg.GetStringSlice("repo.core.meep-swagger-ui.api-bundle")
	//find all the apis and copy them at the location above
	for _, target := range apiBundle {
		apiSrcFile := utils.RepoCfg.GetString("repo.core." + target + ".api")
		if apiSrcFile == "" {
			continue
		}
		apiSrcPath := gitDir + "/" + apiSrcFile
		if apiSrcPath != "" {
			apiDstPath := binDir + "/" + target + "-api.yaml"
			if verbose {
				fmt.Println("    Copying: " + apiSrcPath + " --> " + apiDstPath)
			}

			cmd = exec.Command("cp", apiSrcPath, apiDstPath)
			_, err = utils.ExecuteCmd(cmd, cobraCmd)
			if err != nil {
				fmt.Println("Failed to copy: ", err)
				return
			}

			//find which format style version (standard) it is based on
			nodePort := utils.RepoCfg.GetString("repo.core." + target + ".nodeport")
			switch findFormatStyle(apiDstPath, cobraCmd) {
			case "openapi":
				_ = replaceOpenApiStyle(apiDstPath, nodeIp, nodePort, cobraCmd)
				fmt.Println("Failed to parse/update an openapi yaml file")
				return
			case "swagger":
				err = replaceSwaggerStyle(apiDstPath, nodeIp, nodePort, cobraCmd)
				if err != nil {
					fmt.Println("Failed to parse/update a swagger yaml file")
					return
				}

			default:
			}

			//update the string to update the drop-down menu in the index.html file of /api
			cmd = exec.Command("grep", "title:", apiDstPath)
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
			//urls = urls + `{"name": "` + title + `", "url": "` + target + `-api.yaml"},`
			urls = urls + `{"name": "` + title + `", "url": "http:\/\/` + nodeIp + `:` + ctrlEnginePort + `\/api\/` + target + `-api.yaml"},`
		}
	}

	//update swagger-ui index file
	urls = urls + " ],"
	sedString := "s/" + urlStringToReplace + "/urls: " + urls + "/g"
	cmd = exec.Command("sed", "-i", sedString, binDir+"/index.html")
	_, err = utils.ExecuteCmd(cmd, cobraCmd)
	if err != nil {
		fmt.Println("Failed to sed: ", err)
		return
	}
}

func findFormatStyle(apiPath string, cobraCmd *cobra.Command) string {

	cmd := exec.Command("grep", "openapi: ", apiPath)
	linePresent, _ := utils.ExecuteCmd(cmd, cobraCmd)
	if linePresent != "" {
		//no need to check for openApi version, we handle all the same for now
		return "openapi"
	}

	cmd = exec.Command("grep", "swagger: ", apiPath)
	linePresent, _ = utils.ExecuteCmd(cmd, cobraCmd)
	if linePresent != "" {
		//no need to check for swagger version, we handle all the same for now
		return "swagger"
	}
	return ""
}

func replaceOpenApiStyle(apiPath string, nodeIp string, nodePort string, cobraCmd *cobra.Command) error {
	fmt.Println("No support for openApi files yet!")
	return nil
}

func replaceSwaggerStyle(apiPath string, nodeIp string, nodePort string, cobraCmd *cobra.Command) error {
	//find if host line already exist in the file, if it does, remove it
	cmd := exec.Command("grep", "host: ", apiPath)
	hostLine, _ := utils.ExecuteCmd(cmd, cobraCmd)
	if hostLine != "" {
		hostLine = strings.TrimSpace(hostLine)
		sedHostLine := "/" + hostLine + "/d"
		cmd = exec.Command("sed", "-i", sedHostLine, apiPath)
		_, err := utils.ExecuteCmd(cmd, cobraCmd)
		if err != nil {
			fmt.Println("Failed to sed: ", err)
			return err
		}
	}

	// If there is both a node IP & port - fix basepath so Try It Out works
	if nodeIp != "" && nodePort != "" {
		//find the basepath line in the file and append the host line
		cmd = exec.Command("grep", "basePath: ", apiPath)
		basePath, err := utils.ExecuteCmd(cmd, cobraCmd)
		if err != nil {
			fmt.Println("Failed to grep: ", err)
			return err
		}

		if basePath == "" {
			fmt.Println("Error: basePath shouldn't be empty")
			return err
		}
		newHostLine := "host: " + nodeIp + ":" + nodePort
		newBasePath := strings.Replace(basePath, `/`, `\/`, -1)
		//removing the CR/LF at end of line
		newBasePath = newBasePath[:len(newBasePath)-1]

		sedBasePathLine := "/" + newBasePath + "/a" + newHostLine
		cmd = exec.Command("sed", "-i", sedBasePathLine, apiPath)
		_, err = utils.ExecuteCmd(cmd, cobraCmd)
		if err != nil {
			fmt.Println("Failed to sed: ", err)
			return err
		}
	}

	return nil
}

func buildGoApp(targetName string, cobraCmd *cobra.Command) {
	fmt.Println("--", targetName, "--")
	target := utils.RepoCfg.GetStringMapString("repo.core." + targetName)
	gitDir := viper.GetString("meep.gitdir")
	srcDir := gitDir + "/" + target["src"]
	binDir := gitDir + "/" + target["bin"]
	codecovCapable := utils.RepoCfg.GetBool("repo.core." + targetName + ".codecov")
	lintEnabled := utils.RepoCfg.GetBool("repo.core." + targetName + ".lint")

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
	if lintEnabled && !buildNolint {
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
	fixDeps(targetName, cobraCmd)

	// build
	buildFlags := utils.RepoCfg.GetStringSlice("repo.core." + targetName + ".build-flags")
	if buildCodecov && codecovCapable {
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

	// Copy Dockerfile
	cmd = exec.Command("cp", "Dockerfile", binDir)
	cmd.Dir = srcDir
	out, err = utils.ExecuteCmd(cmd, cobraCmd)
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println(out)
	}
}

func fixDeps(targetName string, cobraCmd *cobra.Command) {
	target := utils.RepoCfg.GetStringMapString("repo.core." + targetName)
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
