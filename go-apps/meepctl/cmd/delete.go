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
  * dep:  Dependency containers
  * all:  All containers
			`,
	Example: `  # Delete all containers
    meepctl delete all
  # Delete only AdvantEDGE core containers
    meepctl delete core
			`,
	Args:      cobra.ExactValidArgs(1),
	ValidArgs: []string{"all", "dep", "core"},
	Run: func(cmd *cobra.Command, args []string) {
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
		utils.InitRepoConfig()
		if group == "all" {
			deleteCore(cmd)
			deleteDep(cmd)
		} else if group == "core" {
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
	k8sDelete("meep-virt-engine", cobraCmd)
	deleteVirtEngine(cobraCmd)
	k8sDelete("meep-webhook", cobraCmd)
	k8sDelete("meep-mg-manager", cobraCmd)
	k8sDelete("meep-tc-engine", cobraCmd)
	k8sDelete("meep-mon-engine", cobraCmd)
	k8sDelete("meep-ctrl-engine", cobraCmd)
	deleteMeepUserAccount(cobraCmd)
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
	gitdir := viper.GetString("meep.gitdir") + "/"

	// Local storage bindings
	// NOTE: Helm charts don't remove pvc for statefulsets because helm did not create them
	// Run in separate threads in order to complete uninstall successfully
	messages := make(chan string)
	go k8sDeletePvc("data-meep-elasticsearch-data-0", cobraCmd, messages)
	go k8sDeletePvc("data-meep-elasticsearch-master-0", cobraCmd, messages)
	go k8sDeletePvc("data-meep-elasticsearch-master-1", cobraCmd, messages)

	k8sDelete("meep-redis", cobraCmd)
	k8sDelete("meep-kube-state-metrics", cobraCmd)
	k8sDelete("meep-metricbeat", cobraCmd)
	k8sDelete("meep-couchdb", cobraCmd)
	k8sDelete("meep-kibana", cobraCmd)
	k8sDelete("meep-filebeat", cobraCmd)
	k8sDelete("meep-curator", cobraCmd)
	k8sDelete("meep-elasticsearch", cobraCmd)
	k8sDelete("meep-docker-registry", cobraCmd)

	// Wait for all pvc delete routines to complete
	// NOTE: Must be checked after deleting elastic
	for i := 0; i < 3; i++ {
		fmt.Println(<-messages)
	}

	// Local storage bindings
	// @TODO move to respective charts
	cmd := exec.Command("kubectl", "delete", "-f", gitdir+utils.RepoCfg.GetString("repo.dep.elastic.es.pv"))
	out, err := utils.ExecuteCmd(cmd, cobraCmd)
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println(out)
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

func k8sDeletePvc(pvc string, cobraCmd *cobra.Command, messages chan string) {
	cmd := exec.Command("kubectl", "delete", "pvc", pvc)
	out, err := utils.ExecuteCmd(cmd, cobraCmd)
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println(out)
	}
	messages <- "Deleted pvc: " + pvc
}
