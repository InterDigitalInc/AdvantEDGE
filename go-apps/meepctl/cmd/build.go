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

var buildCodecov bool

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build <target>",
	Short: "Build core components",
	Long: `Build core components

AdvantEDGE is composed of a collection of micro-services.

Build command genrates AdvantEDGE binaries.

Valid targets:`,
	Example: `# Build all components
	meepctl build all
# Build meep-ctrl-engine component only
	meepctl build meep-ctrl-engine
		`,
	Args: cobra.ExactValidArgs(1),
	// WARNING -- meep-frontend comes before meep-ctrl-engine so that "all" works
	ValidArgs: []string{"all", "meep-frontend", "meep-ctrl-engine", "meep-initializer", "meep-mg-manager", "meep-mon-engine", "meep-tc-engine", "meep-tc-sidecar", "meep-virt-engine"},

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
			fmt.Println("[flag] codecov:", buildCodecov)
			fmt.Println("[flag] verbose:", v)
			fmt.Println("[flag] time:", t)
		}

		start := time.Now()
		utils.InitRepoConfig()
		if target == "all" {
			buildAll(cmd)
		} else {
			build(target, cmd)
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

	if targetName == "meep-frontend" {
		buildFrontend(targetName, cobraCmd)
	} else {
		buildGoApp(targetName, cobraCmd)
	}
}

func buildFrontend(targetName string, cobraCmd *cobra.Command) {
	fmt.Println("--", targetName, "--")
	target := utils.RepoCfg.GetStringMapString("repo.core." + targetName)
	gitDir := viper.GetString("meep.gitdir")
	srcDir := gitDir + "/" + target["src"]
	binDir := gitDir + "/" + target["bin"]

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

	//build
	fmt.Println("   + building " + targetName)
	cmd = exec.Command("npm", "run", "build")
	cmd.Dir = srcDir
	out, err = utils.ExecuteCmd(cmd, cobraCmd)
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println(out)
	}
}

func buildGoApp(targetName string, cobraCmd *cobra.Command) {
	fmt.Println("--", targetName, "--")
	target := utils.RepoCfg.GetStringMapString("repo.core." + targetName)
	gitDir := viper.GetString("meep.gitdir")
	srcDir := gitDir + "/" + target["src"]
	binDir := gitDir + "/" + target["bin"]
	codecovCapable := utils.RepoCfg.GetBool("repo.core." + targetName + ".codecov")

	// dependencies
	fmt.Println("   + checking external dependencies")
	cmd := exec.Command("go", "mod", "vendor")
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
			if targetName == "meep-ctrl-engine" {
				checkMeepCtrlEngineFeDep(targetName, dep, depDir, cobraCmd)
			}
		}
	}

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
}

func checkMeepCtrlEngineFeDep(targetName string, dep string, depDir string, cobraCmd *cobra.Command) {
	verbose, _ := cobraCmd.Flags().GetBool("verbose")
	target := utils.RepoCfg.GetStringMapString("repo.core." + targetName)
	gitDir := viper.GetString("meep.gitdir")
	srcDir := gitDir + "/" + depDir
	dstDir := gitDir + "/" + target["bin"]

	if dep == "meep-frontend" {
		// Only copy FE if a newer FE exists
		if _, err := os.Stat(srcDir); !os.IsNotExist(err) {
			if verbose {
				fmt.Println("    Using: " + srcDir)
			}
			cmd := exec.Command("rm", "-r", dstDir+"/static")
			utils.ExecuteCmd(cmd, cobraCmd)
			cmd = exec.Command("cp", "-r", srcDir, dstDir)
			utils.ExecuteCmd(cmd, cobraCmd)
			cmd = exec.Command("mv", "dist", "static")
			cmd.Dir = dstDir
			utils.ExecuteCmd(cmd, cobraCmd)
		} else if verbose {
			fmt.Println("    Using: " + dstDir + "/static")
		}
	}
}

func fixDeps(targetName string, cobraCmd *cobra.Command) {
	target := utils.RepoCfg.GetStringMapString("repo.core." + targetName)
	gitDir := viper.GetString("meep.gitdir")
	srcDir := gitDir + "/" + target["src"]

	switch targetName {
	case "meep-initializer", "meep-mon-engine":
		cmd := exec.Command("rm", srcDir+"/vendor/k8s.io/client-go/tools/cache/mutation_cache.go")
		utils.ExecuteCmd(cmd, cobraCmd)
		cmd = exec.Command("rm", "-Rf", srcDir+"/vendor/github.com/hashicorp/golang-lru")
		utils.ExecuteCmd(cmd, cobraCmd)
	}
}
