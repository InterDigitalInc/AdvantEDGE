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
	Use:   "deploy <group>",
	Short: "Deploy containers on the K8s cluster",
	Long: `Deploy containers on the K8s cluster

AdvantEDGE is composed of a collection of micro-services (a.k.a the groups).

Deploy command starts a group of containers the in the K8s cluster.
Optional registry & tag parameters allows to specify a shared registry & tag for core images.
Default registry is configured in ~/.meepctl.yaml.
Defaut tag is: latest

Valid groups:
  * core: AdvantEDGE core containers
  * dep:  Dependency containers`,
	Example: `  # Deploy AdvantEDGE dependencies
  meepctl deploy dep
  # Delete and re-deploy only AdvantEDGE core containers
  meepctl deploy core --force
  # Deploy AdvantEDGE version 1.0.0 from my.registry.com
  meepctl deploy core --registry my.registry.com --tag 1.0.0`,
	Args:      cobra.ExactValidArgs(1),
	ValidArgs: []string{"dep", "core"},
	Run: func(cmd *cobra.Command, args []string) {
		if !utils.ConfigValidate("") {
			fmt.Println("Fix configuration issues")
			return
		}

		group := args[0]

		registry, _ := cmd.Flags().GetString("registry")
		tag, _ := cmd.Flags().GetString("tag")
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
		if registry == "" {
			registry = viper.GetString("meep.registry")
		}
		fmt.Println("Using docker registry:", registry)

		if group == "core" {
			deployCore(cmd, registry, tag)
		} else if group == "dep" {
			deployDep(cmd)
		}
		elapsed := time.Since(start)
		if t {
			fmt.Println("Took ", elapsed.Round(time.Millisecond).String())
		}
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
	deployCmd.Flags().BoolP("force", "f", false, "Deployed components are deleted and deployed")
	deployCmd.Flags().BoolVar(&deployCodecov, "codecov", false, "Use when deploying code coverage binaries (dev. option)")
	deployCmd.Flags().StringP("registry", "r", "", "Override registry from config file")
	deployCmd.Flags().StringP("tag", "", "latest", "Repo tag to use")
}

func ensureCoreStorage(cobraCmd *cobra.Command) {
	workdir := viper.GetString("meep.workdir") + "/"

	// Local storage strucutre
	cmd := exec.Command("mkdir", "-p", workdir)
	cmd.Args = append(cmd.Args, workdir+"certs")

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
	gitdir := viper.GetString("meep.gitdir") + "/"
	workdir := viper.GetString("meep.workdir") + "/"
	nodeIp := viper.GetString("node.ip")

	// Local storage structure
	cmd := exec.Command("mkdir", "-p", workdir)
	cmd.Args = append(cmd.Args, workdir+"couchdb")
	cmd.Args = append(cmd.Args, workdir+"influxdb")
	cmd.Args = append(cmd.Args, workdir+"grafana")
	cmd.Args = append(cmd.Args, workdir+"docker-registry")
	cmd.Args = append(cmd.Args, workdir+"certs")

	_, err := utils.ExecuteCmd(cmd, cobraCmd)
	if err != nil {
		err = errors.New("Error creating path [" + workdir + "]")
		fmt.Println(err)
	}

	// EXCEPTION #1: Update Cluster IP address in Grafana values.yaml
	cmd = exec.Command("mkdir", "-p", workdir+"tmp")
	_, _ = utils.ExecuteCmd(cmd, cobraCmd)
	valuesGrafana := gitdir + utils.RepoCfg.GetString("repo.dep.grafana.chart") + "/values.yaml"
	cmd = exec.Command("cp", valuesGrafana, workdir+"tmp/grafana-values.yaml")
	_, _ = utils.ExecuteCmd(cmd, cobraCmd)
	str := "s/<CLUSTERIP>/" + nodeIp + "/g"
	cmd = exec.Command("sed", "-i", str, workdir+"tmp/grafana-values.yaml")
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
	ip := viper.GetString("node.ip")

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
	repo = "meep-loc-serv"
	chart = gitdir + utils.RepoCfg.GetString("repo.core.meep-loc-serv.chart")
	flags := utils.HelmFlags(nil, "--set", "image.env.rooturl=http://"+ip)
	k8sDeployCore(repo, registry, tag, chart, flags, cobraCmd)
	//---
	repo = "meep-metrics-engine"
	chart = gitdir + utils.RepoCfg.GetString("repo.core.meep-metrics-engine.chart")
	flags = utils.HelmFlags(nil, "--set", "image.env.rooturl=http://"+ip)
	k8sDeployCore(repo, registry, tag, chart, flags, cobraCmd)
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
	flags = utils.HelmFlags(nil, "--set", "sidecar.image.repository="+registry+"meep-tc-sidecar")
	flags = utils.HelmFlags(flags, "--set", "sidecar.image.tag="+tag)
	flags = utils.HelmFlags(flags, "--set", "webhook.cert="+cert)
	flags = utils.HelmFlags(flags, "--set", "webhook.key="+key)
	flags = utils.HelmFlags(flags, "--set", "webhook.cabundle="+cabundle)
	k8sDeployCore(repo, registry, tag, chart, flags, cobraCmd)
	//---
	repo = "meep-virt-engine"
	chart = gitdir + utils.RepoCfg.GetString("repo.core.meep-virt-engine.chart")
	flags = utils.HelmFlags(nil, "--set", "service.ip="+ip)
	k8sDeploy(repo, chart, flags, cobraCmd)
	deployVirtEngineExt(repo, cobraCmd)
}

func deployDep(cobraCmd *cobra.Command) {
	var repo string
	var chart string
	var flags [][]string
	gitdir := viper.GetString("meep.gitdir") + "/"
	workdir := viper.GetString("meep.workdir") + "/"

	// Storage
	ensureDepStorage(cobraCmd)

	// Runtime
	repo = "meep-docker-registry"
	chart = gitdir + utils.RepoCfg.GetString("repo.dep.docker-registry.chart")
	flags = utils.HelmFlags(nil, "--set", "persistence.location="+workdir+"docker-registry/")
	createRegistryCerts(chart, workdir+"certs", cobraCmd)
	k8sDeploy(repo, chart, flags, cobraCmd)
	//---
	repo = "meep-couchdb"
	chart = gitdir + utils.RepoCfg.GetString("repo.dep.couchdb.chart")
	flags = utils.HelmFlags(nil, "--set", "persistentVolume.location="+workdir+"couchdb/")
	k8sDeploy(repo, chart, flags, cobraCmd)
	//---
	repo = "meep-redis"
	chart = gitdir + utils.RepoCfg.GetString("repo.dep.redis.chart")
	flags = nil
	k8sDeploy(repo, chart, flags, cobraCmd)
	//---
	repo = "meep-influxdb"
	chart = gitdir + utils.RepoCfg.GetString("repo.dep.influxdb.chart")
	flags = utils.HelmFlags(nil, "--set", "persistence.location="+workdir+"influxdb/")
	k8sDeploy(repo, chart, flags, cobraCmd)
	//---
	repo = "meep-grafana"
	chart = gitdir + utils.RepoCfg.GetString("repo.dep.grafana.chart")
	flags = utils.HelmFlags(nil, "--set", "persistentVolume.location="+workdir+"grafana/")
	flags = utils.HelmFlags(flags, "--values", workdir+"tmp/grafana-values.yaml")
	k8sDeploy(repo, chart, flags, cobraCmd)
	//---
	repo = "meep-kube-state-metrics"
	chart = gitdir + utils.RepoCfg.GetString("repo.dep.kube-state-metrics.chart")
	flags = nil
	k8sDeploy(repo, chart, flags, cobraCmd)
	//---
	repo = "meep-ingress"
	chart = gitdir + utils.RepoCfg.GetString("repo.dep.nginx-ingress.chart")
	createIngressCerts(chart, workdir+"certs", cobraCmd)
	flags = nil
	httpPort, httpsPort := getPorts()
	if httpPort != "80" {
		flags = utils.HelmFlags(flags, "--set", "controller.hostNetwork=false")
		flags = utils.HelmFlags(flags, "--set", "controller.dnsPolicy=ClusterFirst")
		flags = utils.HelmFlags(flags, "--set", "controller.daemonset.useHostPort=false")
		flags = utils.HelmFlags(flags, "--set", "controller.service.type=NodePort")
		flags = utils.HelmFlags(flags, "--set", "controller.service.nodePorts.http="+httpPort)
		flags = utils.HelmFlags(flags, "--set", "controller.service.nodePorts.https="+httpsPort)
	}
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
	gitdir := viper.GetString("meep.gitdir") + "/"
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
	virtEngineApp := gitdir + utils.RepoCfg.GetString("repo.core.meep-virt-engine.bin") + "/meep-virt-engine"
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
	cmd := exec.Command("sh", "-c", chart+"/create-k8s-ca-signed-cert.sh --certdir "+certdir)
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

func createRegistryCerts(chart string, certdir string, cobraCmd *cobra.Command) {
	cmd := exec.Command("sh", "-c", chart+"/create-k8s-ca-signed-cert.sh --certdir "+certdir)
	_, _ = utils.ExecuteCmd(cmd, cobraCmd)
}

func createIngressCerts(chart string, certdir string, cobraCmd *cobra.Command) {
	cmd := exec.Command("sh", "-c", chart+"/create-self-signed-cert.sh --certdir "+certdir)
	_, _ = utils.ExecuteCmd(cmd, cobraCmd)
}

func getPorts() (string, string) {
	ports := viper.GetString("meep.ports")
	p := strings.Split(ports, "/")
	return p[0], p[1]
}
