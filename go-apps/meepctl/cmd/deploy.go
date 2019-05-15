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
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/InterDigitalInc/AdvantEDGE/go-apps/meepctl/utils"

	"github.com/roymx/viper"
	"github.com/spf13/cobra"
)

var deployCodecov bool

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy <group> [registry] [tag]",
	Short: "Deploy containers on the K8s cluster",
	Long: `Deploy containers on the K8s cluster

AdvantEDGE is composed of a collection of micro-services (a.k.a the groups).

Deploy command starts a group of containers the in the K8s cluster.
Optional registry & tag parameters allows to specify a shared registry & tag for core images.
Default registry/tag are: local registry & latest

Valid groups:
  * core: AdvantEDGE core containers
  * dep:  Dependency containers
  * all:  All containers
		`,
	Example: `  # Deploy all containers
    meepctl deploy all
  # Delete and re-deploy only AdvantEDGE core containers
    meepctl deploy core --force
  # Deploy AdvantEDGE version 1.0.0 from my.registry.com
	  meepctl deploy core my.registry.com 1.0.0
			`,
	Args: cobra.RangeArgs(1, 3),
	Run: func(cmd *cobra.Command, args []string) {
		group := args[0]
		registry := ""
		if len(args) > 1 {
			registry = args[1]
		}
		tag := "latest"
		if len(args) > 2 {
			tag = args[2]
		}

		f, _ := cmd.Flags().GetBool("force")
		v, _ := cmd.Flags().GetBool("verbose")
		t, _ := cmd.Flags().GetBool("time")
		if v {
			fmt.Println("Deploy called")
			fmt.Println("[arg]  group:", group)
			fmt.Println("[arg]  registry:", registry)
			fmt.Println("[arg]  tag:", tag)
			fmt.Println("[flag] force:", f)
			fmt.Println("[flag] verbose:", v)
			fmt.Println("[flag] time:", t)
		}

		start := time.Now()
		utils.InitRepoConfig()
		if group == "all" {
			deployDep(cmd)
			deployCore(cmd, registry, tag)
		} else if group == "core" {
			deployCore(cmd, registry, tag)
		} else if group == "dep" {
			deployDep(cmd)
		} else {
			fmt.Println("Invalid group ", group)
			fmt.Println(cmd.Long)
		}
		elapsed := time.Since(start)
		if t {
			fmt.Println("Took ", elapsed.Round(time.Millisecond).String())
		}
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deployCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	deployCmd.Flags().BoolP("force", "f", false, "Deployed components are deleted and deployed")
	deployCmd.Flags().BoolVar(&deployCodecov, "codecov", false, "Use when deploying code coverage binaries (dev. option)")
}

func ensureCoreStorage(cobraCmd *cobra.Command) {
	workdir := viper.GetString("meep.workdir") + "/"

	// Local storage strucutre
	cmd := exec.Command("mkdir", "-p", workdir)
	_, err := utils.ExecuteCmd(cmd, cobraCmd)
	if err != nil {
		err = errors.New("Error creating path [" + workdir + "]")
		fmt.Println(err)
	}

	//templates
	templatedir := viper.GetString("meep.gitdir") + "/" + utils.RepoCfg.GetString("repo.core.meep-virt-engine.template")
	cmd = exec.Command("rm", "-rf", workdir+"template-bak")
	_, _ = utils.ExecuteCmd(cmd, cobraCmd)
	cmd = exec.Command("mv", workdir+"template", workdir+"template-bak")
	_, _ = utils.ExecuteCmd(cmd, cobraCmd)
	cmd = exec.Command("cp", "-r", templatedir, workdir+"template")
	_, _ = utils.ExecuteCmd(cmd, cobraCmd)
	//codecov
	cmd = exec.Command("rm", "-rf", workdir+"codecov-bak")
	_, _ = utils.ExecuteCmd(cmd, cobraCmd)
	cmd = exec.Command("mv", workdir+"codecov", workdir+"codecov-bak")
	_, _ = utils.ExecuteCmd(cmd, cobraCmd)
	cmd = exec.Command("mkdir", "-p", workdir+"codecov")
	_, _ = utils.ExecuteCmd(cmd, cobraCmd)
	for _, targetName := range buildCmd.ValidArgs {
		if targetName == "all" {
			continue
		}
		codecovCapable := utils.RepoCfg.GetBool("repo.core." + targetName + ".codecov")
		if codecovCapable {
			cmd = exec.Command("mkdir", "-p", workdir+"codecov/"+targetName)
			_, _ = utils.ExecuteCmd(cmd, cobraCmd)
		}
	}
	//certs
	cmd = exec.Command("mkdir", "-p", workdir+"certs")
	_, _ = utils.ExecuteCmd(cmd, cobraCmd)
}

func ensureDepStorage(cobraCmd *cobra.Command) {
	workdir := viper.GetString("meep.workdir") + "/"

	// Local storage structure
	cmd := exec.Command("mkdir", "-p", workdir)
	cmd.Args = append(cmd.Args, workdir+"couchdb")
	cmd.Args = append(cmd.Args, workdir+"es-data")
	cmd.Args = append(cmd.Args, workdir+"es-master-0")
	cmd.Args = append(cmd.Args, workdir+"es-master-1")
	cmd.Args = append(cmd.Args, workdir+"kibana")

	_, err := utils.ExecuteCmd(cmd, cobraCmd)
	if err != nil {
		err = errors.New("Error creating path [" + workdir + "]")
		fmt.Println(err)
	}

	//copy the yaml files in workdir and apply a modification to the tmp file, original is untouched
	cmd = exec.Command("mkdir", "-p", workdir+"tmp")
	_, _ = utils.ExecuteCmd(cmd, cobraCmd)
	pvES := viper.GetString("meep.gitdir") + "/" + utils.RepoCfg.GetString("repo.dep.elastic.es.pv")
	cmd = exec.Command("cp", pvES, workdir+"tmp/meep-pv-es.yaml")
	_, _ = utils.ExecuteCmd(cmd, cobraCmd)
	valuesFB := viper.GetString("meep.gitdir") + "/" + utils.RepoCfg.GetString("repo.dep.elastic.filebeat.values")
	cmd = exec.Command("cp", valuesFB, workdir+"tmp/filebeat-values.yaml")
	_, _ = utils.ExecuteCmd(cmd, cobraCmd)
	//search and replace in yaml fil
	tmpStr := strings.Replace(workdir, "/", "\\/", -1)
	str := "s/<WORKDIR>/" + tmpStr + "/g"
	cmd = exec.Command("sed", "-i", str, workdir+"tmp/meep-pv-es.yaml")
	_, _ = utils.ExecuteCmd(cmd, cobraCmd)
	cmd = exec.Command("sed", "-i", str, workdir+"tmp/filebeat-values.yaml")
	_, _ = utils.ExecuteCmd(cmd, cobraCmd)

	// Local storage bindings
	// @TODO move to respective chart
	cmd = exec.Command("kubectl", "apply", "-f", workdir+"tmp/meep-pv-es.yaml")
	_, _ = utils.ExecuteCmd(cmd, cobraCmd)
}

func deployCore(cobraCmd *cobra.Command, registry string, tag string) {
	// Storage
	ensureCoreStorage(cobraCmd)
	// Runtime
	if registry != "" {
		registry = registry + "/"
	}
	gitdir := viper.GetString("meep.gitdir") + "/"
	workdir := viper.GetString("meep.workdir") + "/"

	deployMeepUserAccount(cobraCmd)

	//---
	repo := "meep-ctrl-engine"
	chart := gitdir + utils.RepoCfg.GetString("repo.core.meep-ctrl-engine.chart")
	k8sDeployCore(repo, registry, tag, chart, nil, cobraCmd)
	//---
	repo = "meep-mon-engine"
	chart = gitdir + utils.RepoCfg.GetString("repo.core.meep-mon-engine.chart")
	k8sDeployCore(repo, registry, tag, chart, nil, cobraCmd)
	//---
	repo = "meep-tc-engine"
	chart = gitdir + utils.RepoCfg.GetString("repo.core.meep-tc-engine.chart")
	k8sDeployCore(repo, registry, tag, chart, nil, cobraCmd)
	//---
	repo = "meep-mg-manager"
	chart = gitdir + utils.RepoCfg.GetString("repo.core.meep-mg-manager.chart")
	k8sDeployCore(repo, registry, tag, chart, nil, cobraCmd)
	//---
	repo = "meep-webhook"
	chart = gitdir + utils.RepoCfg.GetString("repo.core.meep-webhook.chart")
	cert, key, cabundle := createWebhookCerts(chart, workdir+"certs", cobraCmd)
	flags := utils.HelmFlags(nil, "--set", "sidecar.image.repository="+registry+"meep-tc-sidecar")
	flags = utils.HelmFlags(flags, "--set", "sidecar.image.tag="+tag)
	flags = utils.HelmFlags(flags, "--set", "webhook.cert="+cert)
	flags = utils.HelmFlags(flags, "--set", "webhook.key="+key)
	flags = utils.HelmFlags(flags, "--set", "webhook.cabundle="+cabundle)
	k8sDeployCore(repo, registry, tag, chart, flags, cobraCmd)
	//---
	repo = "meep-virt-engine"
	chart = gitdir + utils.RepoCfg.GetString("repo.core.meep-virt-engine.chart")
	flags = utils.HelmFlags(nil, "--set", "service.ip="+viper.GetString("node.ip"))
	k8sDeploy(repo, chart, flags, cobraCmd)
	deployVirtEngineExt(repo, cobraCmd)
}

func deployDep(cobraCmd *cobra.Command) {
	gitDir := viper.GetString("meep.gitdir") + "/"
	workdir := viper.GetString("meep.workdir") + "/"

	// Storage
	ensureDepStorage(cobraCmd)
	// Runtime
	repo := "meep-elasticsearch"
	chart := utils.RepoCfg.GetString("repo.dep.elastic.es.chart")
	flags := utils.HelmFlags(nil, "--version", utils.RepoCfg.GetString("repo.dep.elastic.es.version"))
	flags = utils.HelmFlags(flags, "--values", gitDir+utils.RepoCfg.GetString("repo.dep.elastic.es.values"))
	k8sDeploy(repo, chart, flags, cobraCmd)
	//---
	repo = "meep-curator"
	chart = gitDir + utils.RepoCfg.GetString("repo.dep.elastic.es-curator.chart")
	flags = nil
	k8sDeploy(repo, chart, flags, cobraCmd)
	//---
	repo = "meep-kibana"
	chart = gitDir + utils.RepoCfg.GetString("repo.dep.elastic.kibana.chart")
	flags = utils.HelmFlags(nil, "--set", "persistentVolume.location="+workdir+"kibana/")
	k8sDeploy(repo, chart, flags, cobraCmd)
	//---
	// Value file is modified, use the tmp/ version
	repo = "meep-filebeat"
	chart = utils.RepoCfg.GetString("repo.dep.elastic.filebeat.chart")
	flags = utils.HelmFlags(nil, "--version", utils.RepoCfg.GetString("repo.dep.elastic.filebeat.version"))
	flags = utils.HelmFlags(flags, "--values", workdir+"tmp/filebeat-values.yaml")
	k8sDeploy(repo, chart, flags, cobraCmd)
	//---
	repo = "meep-couchdb"
	chart = gitDir + utils.RepoCfg.GetString("repo.dep.couchdb.chart")
	flags = utils.HelmFlags(nil, "--set", "persistentVolume.location="+workdir+"couchdb/")
	k8sDeploy(repo, chart, flags, cobraCmd)
	//---
	repo = "meep-redis"
	chart = utils.RepoCfg.GetString("repo.dep.redis.chart")
	flags = utils.HelmFlags(nil, "--version", utils.RepoCfg.GetString("repo.dep.redis.version"))
	flags = utils.HelmFlags(flags, "--values", gitDir+utils.RepoCfg.GetString("repo.dep.redis.values"))
	k8sDeploy(repo, chart, flags, cobraCmd)
	//---
	repo = "meep-kube-state-metrics"
	chart = gitDir + utils.RepoCfg.GetString("repo.dep.k8s.kube-state-metrics.chart")
	flags = nil
	k8sDeploy(repo, chart, flags, cobraCmd)
	//---
	repo = "meep-metricbeat"
	chart = gitDir + utils.RepoCfg.GetString("repo.dep.elastic.metricbeat.chart")
	flags = utils.HelmFlags(nil, "--set", "image.pullPolicy=IfNotPresent")
	k8sDeploy(repo, chart, flags, cobraCmd)
}

func k8sDeployCore(repo string, registry string, tag string, chart string, flags [][]string, cobraCmd *cobra.Command) {
	coreFlags := utils.HelmFlags(flags, "--set", "image.repository="+registry+repo)
	coreFlags = utils.HelmFlags(coreFlags, "--set", "image.tag="+tag)
	codecovCapable := utils.RepoCfg.GetBool("repo.core." + repo + ".codecov")
	if deployCodecov && codecovCapable {
		coreFlags = utils.HelmFlags(coreFlags, "--set", "codecov.enabled=true")
	}
	k8sDeploy(repo, chart, coreFlags, cobraCmd)
}

func k8sDeploy(component string, chart string, flags [][]string, cobraCmd *cobra.Command) {
	force, _ := cobraCmd.Flags().GetBool("force")

	// If release exist && --force, delete
	exist, _ := utils.IsHelmRelease(component, cobraCmd)
	if exist {
		if force {
			_ = utils.HelmDelete(component, cobraCmd)
		} else {
			fmt.Println("Skipping " + component + ": already deployed -- use [-f, --force] flag to force deployment")
			return
		}
	}

	// Deploy
	_ = utils.HelmInstall(component, chart, flags, cobraCmd)
}

func deployVirtEngineExt(component string, cobraCmd *cobra.Command) {
	verbose, _ := cobraCmd.Flags().GetBool("verbose")
	force, _ := cobraCmd.Flags().GetBool("force")
	gitDir := viper.GetString("meep.gitdir") + "/"
	workdir := viper.GetString("meep.workdir") + "/"
	start := time.Now()

	// If release exist && --force, delete
	pid, err := utils.GetProcess(component, cobraCmd)
	if err == nil && pid != "" {
		if force {
			deleteVirtEngine(cobraCmd)
		} else {
			fmt.Println("Skipping " + component + " (ext.): already deployed -- use [-f, --force] flag to force deployment")
			return
		}
	}

	// Deploy
	// ensure directory
	logdir := workdir + "log"
	cmd := exec.Command("mkdir", "-p", logdir)
	_, _ = utils.ExecuteCmd(cmd, cobraCmd)
	// start ext. component
	file, err := os.Create(logdir + "/virt-engine.log")
	if err != nil {
		fmt.Println("Error starting virt.engine (ext.)")
		fmt.Println(err)
		return
	}

	codecovCapable := utils.RepoCfg.GetBool("repo.core." + component + ".codecov")
	virtEngineApp := gitDir + utils.RepoCfg.GetString("repo.core.meep-virt-engine.bin") + "/meep-virt-engine"
	if deployCodecov && codecovCapable {
		codecovFile := workdir + "/codecov/" + component + "/codecov-meep-virt-engine.out"
		_, _ = utils.ExecuteCmd(cmd, cobraCmd)
		cmd = exec.Command(virtEngineApp, "-test.coverprofile="+codecovFile, "__DEVEL--code-cov")
	} else {
		cmd = exec.Command(virtEngineApp)
	}
	cmd.Stdout = file
	cmd.Stderr = file
	if verbose {
		fmt.Println("Args:", cmd.Args)
	}
	err = cmd.Start()
	elapsed := time.Since(start)

	if err != nil {
		fmt.Println("Error starting virt.engine (ext.)")
		fmt.Println(err)
	} else {
		r := utils.FormatResult("Deployed meep-virt-engine (ext.)", elapsed, cobraCmd)
		fmt.Println(r)
	}
}

func deployMeepUserAccount(cobraCmd *cobra.Command) {
	gitdir := viper.GetString("meep.gitdir")

	cmd := exec.Command("kubectl", "create", "-f", gitdir+"/"+utils.RepoCfg.GetString("repo.core.meep-user.service-account"))
	_, _ = utils.ExecuteCmd(cmd, cobraCmd)

	cmd = exec.Command("kubectl", "create", "-f", gitdir+"/"+utils.RepoCfg.GetString("repo.core.meep-user.cluster-role-binding"))
	_, _ = utils.ExecuteCmd(cmd, cobraCmd)
}

func createWebhookCerts(chart string, certdir string, cobraCmd *cobra.Command) (string, string, string) {
	cmd := exec.Command("sh", "-c", chart+"/webhook-create-signed-cert.sh --certdir "+certdir)
	_, _ = utils.ExecuteCmd(cmd, cobraCmd)
	cmd = exec.Command("sh", "-c", "cat "+certdir+"/server-cert.pem | base64 -w0")
	cert, _ := utils.ExecuteCmd(cmd, cobraCmd)
	cmd = exec.Command("sh", "-c", "cat "+certdir+"/server-key.pem | base64 -w0")
	key, _ := utils.ExecuteCmd(cmd, cobraCmd)
	cmd = exec.Command("kubectl", "config", "view", "--raw", "--minify", "--flatten",
		"-o=jsonpath='{.clusters[].cluster.certificate-authority-data}'")
	cabundle, _ := utils.ExecuteCmd(cmd, cobraCmd)
	return cert, key, cabundle
}
