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
	"strings"
	"time"

	"github.com/InterDigitalInc/AdvantEDGE/go-apps/meepctl/utils"

	"github.com/spf13/cobra"
)

const meepctlVersion = "1.4.2"
const na = "NA"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version <group>",
	Short: "Display version information",
	Long: `Display version information

AdvantEDGE is composed of a collection of components running as micro-services/applications.

Versions command collects and displays version of core & dependency components

Valid groups:
  * core: AdvantEDGE core containers
  * dep:  Dependency applications
  * all:  All containers and applications
  * <none>: Displays the version of the meepctl tool
                        `,
	Example: `  # Displays Versions of all containers
  meepctl version all
  # Display versions of only AdvantEDGE core containers
  meepctl version core
                        `,
	Args:      cobra.MaximumNArgs(1),
	ValidArgs: []string{"all", "dep", "core"},
	Run: func(cmd *cobra.Command, args []string) {
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
		ver := formatVersion("meepctl", meepctlVersion, "", "")
		fmt.Println(ver)
		repoVer := formatVersion(".meepctl-repocfg.yaml", utils.RepoCfg.GetString("version"), "", "")
		fmt.Println(repoVer)
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
	},
}

type versionInfo struct {
	Name      string `json:"name"`
	Version   string `json:"version,omitempty"`
	VersionID string `json:"id,omitempty"`
	BuildID   string `json:"build,omitempty"`
}

var corePodsNameMap = []string{
	"meep-ctrl-engine",
	"meep-webhook",
	"meep-mg-manager",
	"meep-mon-engine",
	"meep-redis",
	"meep-tc-engine",
	"meep-rnis",
	"meep-loc-serv",
	"meep-influxdb",
	"grafana",
	"couchdb",
	"kube-state-metrics",
	"docker-registry",
}

var depPodsNameMap = []string{"weave"}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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
	outVer := getPodVersions(depPodsNameMap, cobraCmd)
	for i := range depPodsNameMap {
		if p, ok := outVer[depPodsNameMap[i]]; ok {
			fmt.Println(formatVersion(p.Name, p.Version, p.VersionID, ""))
		} else {
			fmt.Println(formatVersion(depPodsNameMap[i], na, "", ""))
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
			if vi.Name != "" && contains(podList, vi.Name) {
				outMap[vi.Name] = vi

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

func versionsDep(cobraCmd *cobra.Command) {
	getHelmVersion(cobraCmd)
	getDockerVersion(cobraCmd)
	getKubernetesVersion(cobraCmd)
}

func versionsCore(cobraCmd *cobra.Command) {
	outVer := getPodVersions(corePodsNameMap, cobraCmd)
	for i := range corePodsNameMap {
		if p, ok := outVer[corePodsNameMap[i]]; ok {
			fmt.Println(formatVersion(p.Name, p.Version, p.VersionID, p.BuildID))
		} else {
			fmt.Println(formatVersion(corePodsNameMap[i], na, "", ""))
		}
	}
}
