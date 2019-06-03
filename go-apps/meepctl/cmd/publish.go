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
	"time"

	"github.com/InterDigitalInc/AdvantEDGE/go-apps/meepctl/utils"

	"github.com/spf13/cobra"
)

// publishCmd represents the publish command
var publishCmd = &cobra.Command{
	Use:   "publish <registry> [tag]",
	Short: "Publish image on a registry",
	Long: `Publish image on a registry

AdvantEDGE is composed of a collection of micro-services.

Publish command is a utility function that pushes AdvatEDGE core images to a specified registry.
Core images are pushed from local-registry/core-repo:latest to registry/core-repo:tag (no tag = latest)
`,
	Example: `  # Publish images to your.registry.com:latest
    meepctl publish your.registry.com
  # Publish images to your.registry.com:latest:1.1.0
    meepctl publish your.registry.com 1.1.0`,
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		registry := args[0]
		tag := "latest"
		if len(args) > 1 {
			tag = args[1]
		}

		v, _ := cmd.Flags().GetBool("verbose")
		t, _ := cmd.Flags().GetBool("time")
		if v {
			fmt.Println("Publish called")
			fmt.Println("[arg]  registry:", registry)
			fmt.Println("[arg]  tag:", tag)
			fmt.Println("[flag] verbose:", v)
			fmt.Println("[flag] time:", t)
		}

		start := time.Now()
		publishImage(cmd, registry, "meep-ctrl-engine", tag)
		publishImage(cmd, registry, "meep-mon-engine", tag)
		publishImage(cmd, registry, "meep-loc-serv", tag)
		publishImage(cmd, registry, "meep-tc-engine", tag)
		publishImage(cmd, registry, "meep-tc-sidecar", tag)
		publishImage(cmd, registry, "meep-mg-manager", tag)
		publishImage(cmd, registry, "meep-webhook", tag)
		elapsed := time.Since(start)
		if t {
			fmt.Println("Took ", elapsed.Round(time.Millisecond).String())
		}

	},
}

func init() {
	rootCmd.AddCommand(publishCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// publishCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// publishCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
func publishImage(cobraCmd *cobra.Command, registry string, repo string, tag string) {
	// docker tag meep-ctrl-engine:latest dev.advantedge.com/meep-ctrl-engine:0.2.17
	// docker push dev.advantedge.com/meep-ctrl-engine:0.2.17

	exist, _ := utils.IsDockerImage(repo, cobraCmd)
	if exist {
		localTag, _ := utils.GetDockerTag(repo, cobraCmd)
		localRepo := repo + ":" + localTag
		newRepo := registry + "/" + repo + ":" + tag
		_ = utils.TagDockerImage(localRepo, newRepo, cobraCmd)
		_ = utils.PushDockerImage(newRepo, cobraCmd)
	} else {
		fmt.Println("Image", repo, ":latest not found")
	}
}
