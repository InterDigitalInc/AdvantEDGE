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
	"time"

	"github.com/InterDigitalInc/AdvantEDGE/go-apps/meepctl/utils"

	"github.com/spf13/cobra"
)

type DeleteData struct {
	coreApps []string
	depApps  []string
}

const deleteDesc = `Delete containers from the K8s cluster

AdvantEDGE is composed of a collection of micro-services (a.k.a the groups).

Delete command removes a group of containers from the K8s cluster.`

const deleteExample = `  # Delete dependency containers
  meepctl delete dep
  # Delete only AdvantEDGE core containers
  meepctl delete core`

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:       "delete <group>",
	Short:     "Delete containers from the K8s cluster",
	Long:      deleteDesc,
	Example:   deleteExample,
	Args:      cobra.ExactValidArgs(1),
	ValidArgs: nil,
	Run:       deleteRun,
}

var deleteData DeleteData

func init() {
	// Get targets from repo config file
	deleteData.coreApps = utils.GetTargets("repo.core.go-apps")
	deleteData.depApps = utils.GetTargets("repo.dep")

	// Configure the list of valid arguments
	deleteCmd.ValidArgs = []string{"dep", "core"}

	// Add list of arguments to Example usage
	deleteCmd.Example += "\n\nValid Targets:"
	for _, arg := range deleteCmd.ValidArgs {
		deleteCmd.Example += "\n  * " + arg
	}

	// Add command
	rootCmd.AddCommand(deleteCmd)
}

func deleteRun(cmd *cobra.Command, args []string) {
	if !utils.ConfigValidate("") {
		fmt.Println("Fix configuration issues")
		return
	}

	group := args[0]
	v, _ := cmd.Flags().GetBool("verbose")
	t, _ := cmd.Flags().GetBool("time")
	if v {
		fmt.Println("Delete called")
		fmt.Println("[arg]  group:", group)
		fmt.Println("[flag] verbose:", v)
		fmt.Println("[flag] time:", t)
	}

	start := time.Now()
	if group == "core" {
		deleteApps(deleteData.coreApps, cmd)
	} else if group == "dep" {
		deleteApps(deleteData.depApps, cmd)
	}

	elapsed := time.Since(start)
	if t {
		fmt.Println("Took ", elapsed.Round(time.Millisecond).String())
	}
}

func deleteApps(apps []string, cobraCmd *cobra.Command) {
	messages := make(chan string)
	for _, app := range apps {
		go k8sDelete(app, cobraCmd, messages)
	}

	for i := 0; i < len(apps); i++ {
		fmt.Println(<-messages)
	}
}

func k8sDelete(component string, cobraCmd *cobra.Command, messages chan string) {
	// If release exist
	exist, _ := utils.IsHelmRelease(component, cobraCmd)
	if exist {
		// Delete
		err := utils.HelmDelete(component, cobraCmd)
		if err != nil {
			fmt.Println("Helm delete failed with Error: ", err)
		}
		messages <- "Deleted " + component
	} else {
		messages <- "Missing " + component
	}
}
