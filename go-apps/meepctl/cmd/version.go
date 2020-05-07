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
	"encoding/json"
	"fmt"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/InterDigitalInc/AdvantEDGE/go-apps/meepctl/utils"

	"github.com/spf13/cobra"
)

type VersionData struct {
	coreApps []string
	depApps  []string
}

type versionInfo struct {
	Name      string `json:"name"`
	Version   string `json:"version,omitempty"`
	VersionID string `json:"id,omitempty"`
	BuildID   string `json:"build,omitempty"`
}

const meepctlVersion = "1.4.3"
const na = "NA"

const versionDesc = `Display version information

AdvantEDGE is composed of a collection of components running as micro-services/applications.

Versions command collects and displays version of core & dependency components

Valid groups:
  * core: AdvantEDGE core containers
  * dep:  Dependency applications
  * all:  All containers and applications
  * <none>: Displays the version of the meepctl tool`

const versionExample = `  # Displays Versions of all containers
  meepctl version all
  # Display versions of only AdvantEDGE core containers
  meepctl version core`

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:       "version <group>",
	Short:     "Display version information",
	Long:      versionDesc,
	Example:   versionExample,
	Args:      cobra.MaximumNArgs(1),
	ValidArgs: nil,
	Run:       versionRun,
}

var versionData VersionData

func init() {
	// Get targets from repo config file
	versionData.coreApps = utils.GetTargets("repo.core.go-apps")
	// versionData.coreApps = append(versionData.coreApps, utils.GetTargets("repo.sandbox.go-apps")...)
	sort.Strings(versionData.coreApps)

	versionData.depApps = utils.GetTargets("repo.dep")
	sort.Strings(versionData.depApps)

	// Configure the list of valid arguments
	versionCmd.ValidArgs = []string{"all", "core", "dep"}

	// Add list of arguments to Example usage
	versionCmd.Example += "\n\nValid Targets:"
	for _, arg := range versionCmd.ValidArgs {
		versionCmd.Example += "\n  * " + arg
	}

	// Add command
	rootCmd.AddCommand(versionCmd)
}

func versionRun(cmd *cobra.Command, args []string) {
	if !utils.ConfigValidate("") {
		fmt.Println("Fix configuration issues")
		return
	}

	group := ""
	if len(args) > 0 {
		group = args[0]
	}
	v, _ := cmd.Flags().GetBool("verbose")
	t, _ := cmd.Flags().GetBool("time")
	if v {
		fmt.Println("Version called")
		fmt.Println("[arg]  group:", group)
		fmt.Println("[flag] verbose:", v)
		fmt.Println("[flag] time:", t)
	}

	start := time.Now()

	// Print meepctl version
	ver := formatVersion("meepctl", meepctlVersion, "", "")
	fmt.Println(ver)
	repoVer := formatVersion(".meepctl-repocfg.yaml", utils.RepoCfg.GetString("version"), "", "")
	fmt.Println(repoVer)

	// Print
	if group == "all" {
		versionsDep(cmd)
		versionsCore(cmd)
	} else if group == "core" {
		versionsCore(cmd)
	} else if group == "dep" {
		versionsDep(cmd)
	}

	elapsed := time.Since(start)
	if t {
		fmt.Println("Took ", elapsed.Round(time.Millisecond).String())
	}
}

func formatVersion(name string, version string, versionID string, buildID string) string {
	var info versionInfo
	info.Name = name
	info.Version = version
	info.VersionID = versionID
	info.BuildID = buildID
	v, _ := json.MarshalIndent(info, "", "  ")
	return string(v)
}

func versionsCore(cobraCmd *cobra.Command) {
	// Get core versions
	outVer := getPodVersions(versionData.coreApps, cobraCmd)
	for _, app := range versionData.coreApps {
		if p, ok := outVer[app]; ok {
			fmt.Println(formatVersion(p.Name, p.Version, p.VersionID, p.BuildID))
		} else {
			fmt.Println(formatVersion(app, na, "", ""))
		}
	}
}

func versionsDep(cobraCmd *cobra.Command) {
	// Get dependency versions
	outVer := getPodVersions(versionData.depApps, cobraCmd)
	for _, app := range versionData.depApps {
		if p, ok := outVer[app]; ok {
			fmt.Println(formatVersion(p.Name, p.Version, p.VersionID, p.BuildID))
		} else {
			fmt.Println(formatVersion(app, na, "", ""))
		}
	}

	// Gert additional dependency versions
	getHelmVersion(cobraCmd)
	getDockerVersion(cobraCmd)
	getKubernetesVersion(cobraCmd)
}

func getHelmVersion(cobraCmd *cobra.Command) {
	clientStr := formatVersion("helm client", na, "", "")
	serverStr := formatVersion("helm server", na, "", "")
	cmd := exec.Command("helm", "version")
	output, err := utils.ExecuteCmd(cmd, cobraCmd)
	if err != nil {
		fmt.Println("Error getting helm version\n", err)
	} else {
		output = strings.Replace(output, "\"", "", -1)
		outAll := strings.Split(output, "}")
		outClient := outAll[0]
		outServer := outAll[1]
		//client part
		out := strings.Split(outClient, ",")
		clientStr = formatVersion("helm client", strings.Split(out[0], ":")[2], strings.Split(out[1], ":")[1], "")
		//server part
		out = strings.Split(outServer, ",")
		serverStr = formatVersion("helm server", strings.Split(out[0], ":")[2], strings.Split(out[1], ":")[1], "")
	}

	fmt.Println(clientStr)
	fmt.Println(serverStr)
}

func getDockerVersion(cobraCmd *cobra.Command) {
	clientStr := formatVersion("docker client", na, "", "")
	serverStr := formatVersion("docker server", na, "", "")
	cmd := exec.Command("docker", "version")
	output, err := utils.ExecuteCmd(cmd, cobraCmd)
	if err != nil {
		fmt.Println("Error getting docker version\n", err)
	} else {
		output = strings.Replace(output, " ", "", -1)
		output = strings.Replace(output, "\n", ":", -1)
		out := strings.Split(output, ":")

		clientStr = formatVersion("docker client", out[3], out[9], "")
		serverStr = formatVersion("docker server", out[24], out[30], "")
	}
	fmt.Println(clientStr)
	fmt.Println(serverStr)
}

func getKubernetesVersion(cobraCmd *cobra.Command) {
	clientStr := formatVersion("k8s client", na, "", "")
	serverStr := formatVersion("k8s server", na, "", "")
	cmd := exec.Command("kubectl", "version")
	output, err := utils.ExecuteCmd(cmd, cobraCmd)
	if err != nil {
		fmt.Println("Error getting kubernetes version\n", err)
	} else {
		output = strings.Replace(output, "\"", "", -1)
		output = strings.Replace(output, "\n", ":", -1)
		out := strings.Split(output, ":")

		outVersion := strings.Split(out[4], ",")
		outGitCommit := strings.Split(out[5], ",")
		clientStr = formatVersion("k8s client", outVersion[0], outGitCommit[0], "")

		outVersion = strings.Split(out[17], ",")
		outGitCommit = strings.Split(out[18], ",")
		serverStr = formatVersion("k8s server", outVersion[0], outGitCommit[0], "")
	}
	fmt.Println(clientStr)
	fmt.Println(serverStr)

	/* weave section as part of kubernetes */
	k8sDepPodNames := []string{"weave"}
	outVer := getPodVersions(k8sDepPodNames, cobraCmd)
	for _, podName := range k8sDepPodNames {
		if p, ok := outVer[podName]; ok {
			fmt.Println(formatVersion(p.Name, p.Version, p.VersionID, ""))
		} else {
			fmt.Println(formatVersion(podName, na, "", ""))
		}
	}
}

// contains tells whether a contains x.
func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

/* just a generic function that gets all the pod from all namespaces, filtering should be done by the caller */
func getPodVersions(podList []string, cobraCmd *cobra.Command) map[string]*versionInfo {
	outMap := make(map[string]*versionInfo)
	cmd := exec.Command("kubectl", "get", "pods", "--all-namespaces", "-o", "jsonpath={range .items[*]}{\"{\\\"name\\\":\\\"\"}{.status.containerStatuses[].name}{\"\\\",\"}{\"\\\"version\\\":\\\"\"}{.status.containerStatuses[].image}{\"\\\",\"}{\"\\\"id\\\":\\\"\"}{.status.containerStatuses[].imageID}{\"\\\"}\\n\"}")
	output, err := utils.ExecuteCmd(cmd, cobraCmd)
	if err != nil {
		fmt.Println("Error getting pods version\n", err)
	} else {
		pods := strings.Split(output, "\n")
		for i := range pods {

			vi := new(versionInfo)
			err := json.Unmarshal([]byte(pods[i]), &vi)
			if err != nil {
				continue
			}
			svcName := getSvcName(vi.Name)
			if svcName != "" && contains(podList, svcName) {
				outMap[svcName] = vi

				// Build ID (custom docker image label)
				cmd := exec.Command("docker", "image", "inspect", "--format", "{{ index .Config.Labels \"MeepVersion\"}}", vi.Version)
				vi.BuildID, _ = utils.ExecuteCmd(cmd, cobraCmd)
				vi.BuildID = strings.TrimSuffix(vi.BuildID, "\n")

				// // Image name
				// tv := strings.Split(vi.Version, ":")
				// vi.Version = tv[len(tv)-1]

				// Image ID
				tid := strings.Split(vi.VersionID, ":")
				vi.VersionID = tid[len(tid)-1]
			}
		}
	}
	return outMap
}

func getSvcName(containerName string) string {
	svcName := containerName
	switch containerName {
	case "couchdb":
		svcName = "meep-couchdb"
	case "docker-registry":
		svcName = "meep-docker-registry"
	case "grafana":
		svcName = "meep-grafana"
	case "kube-state-metrics":
		svcName = "meep-kube-state-metrics"
	case "nginx-ingress-controller":
		svcName = "meep-ingress"
	}
	return svcName
}
