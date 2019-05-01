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
	k8sDelete("meep-initializer", cobraCmd)
	k8sDelete("meep-mg-manager", cobraCmd)
	k8sDelete("meep-tc-engine", cobraCmd)
	k8sDelete("meep-mon-engine", cobraCmd)
	k8sDelete("meep-ctrl-engine", cobraCmd)
	deleteMeepUserAccount(cobraCmd)
}

func deleteVirtEngine(cobraCmd *cobra.Command) {
	utils.InterruptProcess("meep-virt-engine", cobraCmd)
}

func deleteMeepUserAccount(cobraCmd *cobra.Command) {
	gitdir := viper.GetString("meep.gitdir")

	cmd := exec.Command("kubectl", "delete", "-f", gitdir+"/"+utils.RepoCfg.GetString("repo.core.meep-user.service-account"))
	utils.ExecuteCmd(cmd, cobraCmd)

	cmd = exec.Command("kubectl", "delete", "-f", gitdir+"/"+utils.RepoCfg.GetString("repo.core.meep-user.cluster-role-binding"))
	utils.ExecuteCmd(cmd, cobraCmd)
}

func deleteDep(cobraCmd *cobra.Command) {
	gitdir := viper.GetString("meep.gitdir") + "/"

	// Local storage bindings
	// NOTE: Helm charts don't remove pvc for statefulsets because helm did not create them
	// Run in separate threads in order to complete uninstall successfully
	messages := make(chan string)
	go func() {
		cmd := exec.Command("kubectl", "delete", "pvc", "database-storage-couchdb-couchdb-0")
		utils.ExecuteCmd(cmd, cobraCmd)
		messages <- "Deleted couchdb pvc"
	}()
	go func() {
		cmd := exec.Command("kubectl", "delete", "pvc", "data-elastic-elasticsearch-data-0")
		utils.ExecuteCmd(cmd, cobraCmd)
		messages <- "Deleted elastic data-0 pvc"
	}()
	go func() {
		cmd := exec.Command("kubectl", "delete", "pvc", "data-elastic-elasticsearch-master-0")
		utils.ExecuteCmd(cmd, cobraCmd)
		messages <- "Deleted elastic master-0 pvc"
	}()
	go func() {
		cmd := exec.Command("kubectl", "delete", "pvc", "data-elastic-elasticsearch-master-1")
		utils.ExecuteCmd(cmd, cobraCmd)
		messages <- "Deleted elastic master-1 pvc"
	}()

	k8sDelete("meep-redis", cobraCmd)
	k8sDelete("kube-state-metrics", cobraCmd)
	k8sDelete("metricbeat", cobraCmd)
	k8sDelete("couchdb", cobraCmd)
	k8sDelete("kibana", cobraCmd)
	k8sDelete("filebeat", cobraCmd)
	k8sDelete("curator", cobraCmd)
	k8sDelete("elastic", cobraCmd)

	// Wait for all pvc delete routines to complete
	// NOTE: Must be checked after deleting couchdb & elastic
	for i := 0; i < 5; i++ {
		fmt.Println(<-messages)
	}

	// Local storage bindings
	// @TODO move to respective charts
	cmd := exec.Command("kubectl", "delete", "-f", gitdir+utils.RepoCfg.GetString("repo.dep.couchdb.pv"))
	utils.ExecuteCmd(cmd, cobraCmd)
	cmd = exec.Command("kubectl", "delete", "-f", gitdir+utils.RepoCfg.GetString("repo.dep.elastic.es.pv"))
	utils.ExecuteCmd(cmd, cobraCmd)
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

func k8sDelete(component string, cobraCmd *cobra.Command) (err error) {
	err = nil

	// If release exist
	exist, _ := utils.IsHelmRelease(component, cobraCmd)
	if exist {
		// Delete
		err = utils.HelmDelete(component, cobraCmd)
	}
	return err
}
