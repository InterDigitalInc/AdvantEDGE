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
	crds     []string
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
	deleteData.crds, _ = utils.GetResourcePrerequisites("repo.resource-prerequisites.crds")
	deleteData.coreApps = utils.GetTargets("repo.core.go-apps", "deploy")
	deleteData.depApps = utils.GetTargets("repo.dep", "deploy")

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
		// Removing CRDs prevents UT execution without deploy dependency pods.
		// For now, meepctl will install CRDs first but will no longer remove them.
		// deleteCRD(deleteData.crds, cmd)
	}

	elapsed := time.Since(start)
	if t {
		fmt.Println("Took ", elapsed.Round(time.Millisecond).String())
	}
}

// func deleteCRD(apps []string, cobraCmd *cobra.Command) {
// 	for _, crd := range apps {
// 		cmd := exec.Command("kubectl", "delete", "crd", crd)
// 		_, err := utils.ExecuteCmd(cmd, cobraCmd)
// 		if err != nil {
// 			err = errors.New("Error deleting CRD, name: [" + crd + "]")
// 			fmt.Println(err)
// 		}
// 	}
// }

func deleteApps(apps []string, cobraCmd *cobra.Command) {
	for _, app := range apps {
		k8sDelete(app, cobraCmd)
	}
}

func k8sDelete(component string, cobraCmd *cobra.Command) {
	// If release exist
	exist, _ := utils.IsHelmRelease(component, cobraCmd)
	if exist {
		// Delete
		err := utils.HelmDelete(component, cobraCmd)
		if err != nil {
			fmt.Println("Helm delete failed with Error: ", err)
		}
	}
}
