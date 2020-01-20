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
	"os/exec"
	"time"

	"github.com/InterDigitalInc/AdvantEDGE/go-apps/meepctl/utils"

	"github.com/roymx/viper"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete <group>",
	Short: "Delete containers from the K8s cluster",
	Long: `Delete containers from the K8s cluster

AdvantEDGE is composed of a collection of micro-services (a.k.a the groups).

Delete command removes a group of containers from the K8s cluster.

Valid groups:
  * core: AdvantEDGE core containers
  * dep:  Dependency containers`,
	Example: `  # Delete dependency containers
  meepctl delete dep
  # Delete only AdvantEDGE core containers
  meepctl delete core`,
	Args:      cobra.ExactValidArgs(1),
	ValidArgs: []string{"dep", "core"},
	Run: func(cmd *cobra.Command, args []string) {
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
			deleteCore(cmd)
		} else if group == "dep" {
			deleteDep(cmd)
		}
		elapsed := time.Since(start)
		if t {
			fmt.Println("Took ", elapsed.Round(time.Millisecond).String())
		}
	},
}

func deleteCore(cobraCmd *cobra.Command) {

	messages := make(chan string)

	go k8sDelete("meep-virt-engine", cobraCmd, messages)
	go k8sDelete("meep-webhook", cobraCmd, messages)
	go k8sDelete("meep-mg-manager", cobraCmd, messages)
	go k8sDelete("meep-tc-engine", cobraCmd, messages)
	go k8sDelete("meep-mon-engine", cobraCmd, messages)
	go k8sDelete("meep-loc-serv", cobraCmd, messages)
	go k8sDelete("meep-metrics-engine", cobraCmd, messages)
	go k8sDelete("meep-ctrl-engine", cobraCmd, messages)
	deleteVirtEngine(cobraCmd)
	deleteMeepUserAccount(cobraCmd)

	for i := 0; i < 8; i++ {
		fmt.Println(<-messages)
	}

}

func deleteVirtEngine(cobraCmd *cobra.Command) {
	pid, err := utils.GetProcess("meep-virt-engine", cobraCmd)
	if err == nil {
		var timeoutMsg string
		start := time.Now()
		// Try interrupting process first
		utils.InterruptProcess(pid, cobraCmd)
		err = utils.WaitProcess(pid, "5", cobraCmd)
		if err != nil {
			// Kill process if it did not exit before timeout
			utils.KillProcess(pid, cobraCmd)
			err = utils.WaitProcess(pid, "5", cobraCmd)
			if err != nil {
				timeoutMsg = " failed with timeout error: " + err.Error()
			}
		}
		elapsed := time.Since(start)
		r := utils.FormatResult("Deleted meep-virt-engine (ext.)"+timeoutMsg, elapsed, cobraCmd)
		fmt.Println(r)
	}
}

func deleteMeepUserAccount(cobraCmd *cobra.Command) {
	gitdir := viper.GetString("meep.gitdir")

	cmd := exec.Command("kubectl", "delete", "-f", gitdir+"/"+utils.RepoCfg.GetString("repo.core.meep-user.service-account"))
	out, err := utils.ExecuteCmd(cmd, cobraCmd)
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println(out)
	}

	cmd = exec.Command("kubectl", "delete", "-f", gitdir+"/"+utils.RepoCfg.GetString("repo.core.meep-user.cluster-role-binding"))
	out, err = utils.ExecuteCmd(cmd, cobraCmd)
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println(out)
	}
}

func deleteDep(cobraCmd *cobra.Command) {
	// Local storage bindings
	// NOTE: Helm charts don't remove pvc for statefulsets because helm did not create them
	// Run in separate threads in order to complete uninstall successfully
	messages := make(chan string)
	go k8sDelete("meep-redis", cobraCmd, messages)
	go k8sDelete("meep-kube-state-metrics", cobraCmd, messages)
	go k8sDelete("meep-couchdb", cobraCmd, messages)
	go k8sDelete("meep-grafana", cobraCmd, messages)
	go k8sDelete("meep-influxdb", cobraCmd, messages)
	go k8sDelete("meep-docker-registry", cobraCmd, messages)

	// Wait for all pvc delete routines to complete
	for i := 0; i < 6; i++ {
		fmt.Println(<-messages)
	}
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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
