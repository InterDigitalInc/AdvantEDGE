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
	"bufio"
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

// configKibana represents the configKibana command
var configKibana = &cobra.Command{
	Use:   "kibana",
	Short: "Configures Kibana (index pattern, saved objects such as dashboards, visualisations, etc.)",
	Long: `Configures Kibana (index pattern, saved objects such as dashboards, visualisations, etc.)
Any Kibana saved object will be overwritten in the process if the object Id are the same
meepctl config kibana.`,
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if !utils.ConfigValidate("") {
			fmt.Println("Fix configuration issues")
			return
		}

		t, _ := cmd.Flags().GetBool("time")
		start := time.Now()
		deployKibanaDashboards(cmd)

		elapsed := time.Since(start)
		if t {
			fmt.Println("Took ", elapsed.Round(time.Millisecond).String())
		}
	},
}

func init() {
	configCmd.AddCommand(configKibana)
}

func deployKibanaDashboards(cobraCmd *cobra.Command) {

	gitDir := viper.GetString("meep.gitdir") + "/"
	workdir := viper.GetString("meep.workdir")

	start := time.Now()

	//make sure kibana is up and ready to receive messages
	//this only happens when kibana is able to connect to elastic search
	//so elastic search must be up and running as well as kibana for the rest api server to be up
	isKibanaUp := false
	kibanaFailedAttempts := 0
	for !isKibanaUp {
		isKibanaUp = uploadDefaultIndex("DUMMY", cobraCmd)
		if !isKibanaUp {
			kibanaFailedAttempts++
			if kibanaFailedAttempts > 3 {
				elapsed := time.Since(start)
				r := utils.FormatResult("Failure during deployment kibana dashboards", elapsed, cobraCmd)
				fmt.Println(r)
				return
			}
		}
	}

	cmd := exec.Command("cp", gitDir+"dashboards/dashboards.conf", workdir+"/tmp/dashboards.conf")
	_, _ = utils.ExecuteCmd(cmd, cobraCmd)
	//search and replace in yaml file
	tmpStr := strings.Replace(gitDir, "/", "\\/", -1)
	str := "s/<GITDIR>/" + tmpStr + "/g"
	cmd = exec.Command("sed", "-i", str, workdir+"/tmp/dashboards.conf")
	_, _ = utils.ExecuteCmd(cmd, cobraCmd)

	f, _ := os.Open(workdir + "/tmp/dashboards.conf")
	//read file line by line
	line := bufio.NewScanner(f)
	for line.Scan() {
		//dashboard[0] = name, dashboard[1] = location
		dashboard := strings.Split(line.Text(), ":")
		if dashboard != nil && dashboard[0] != "" && dashboard[0][0] != '#' {
			//defaultIndex is reserved
			if dashboard[0] == "defaultIndex" {
				defaultIndex := strings.TrimSpace(dashboard[1])
				_ = uploadDefaultIndex(defaultIndex, cobraCmd)
			} else {
				if len(dashboard) >= 2 {
					dashboard_location := strings.TrimSpace(dashboard[1])
					if dashboard_location[:4] == "http" {
						//location is a http location, it had an extra":", so put back the string together
						dashboard_location = strings.TrimSpace(dashboard[2][2:])
						uploadDashboardHttp(dashboard_location, cobraCmd)
					} else {
						uploadDashboardFile(dashboard_location, cobraCmd)
					}
				}
			}
		}
	}

	err := line.Err()
	if err != nil {
		fmt.Println(err)
	}
	elapsed := time.Since(start)
	r := utils.FormatResult("Deployed kibana dashboards", elapsed, cobraCmd)
	fmt.Println(r)

}

//communicating with Kibana, return true if the following conditions are met:
//- all elastic search(ES) pods are up
//- kibana pod is up
//- kibana connected successfully to ES
func uploadDefaultIndex(indexId string, cobraCmd *cobra.Command) bool {
	kibanaHost := viper.GetString("node.ip")
	verbose, _ := cobraCmd.Flags().GetBool("verbose")

	cmd := exec.Command("curl", "-vX", "POST", "http://"+kibanaHost+":32003/api/kibana/settings/defaultIndex", "-H", "Content-Type: application/json", "-H", "kbn-xsrf: true", "-d", "{\"value\": \""+indexId+"\"}")
	out, err := utils.ExecuteCmd(cmd, cobraCmd)
	if err != nil {
		err = errors.New("Error sending a curl command")
	} else {
		str := string(out)
		isServiceUnavailable := strings.Contains(str, "Service Unavailable")

		if isServiceUnavailable {
			err = errors.New("Error: Service Unavailable")
		}
	}

	if err != nil {
		if verbose {
			fmt.Println("Failed to upload a default index error: " + err.Error())
		}
		return false
	} else {
		return true
	}
}

func uploadDashboardHttp(location string, cobraCmd *cobra.Command) {
	//no support yet for url in kibana 6.4.2... but we can get the file to /tmp and then download the file
	strArray := strings.Split(location, "/")
	tmpLocation := "/tmp/" + strArray[len(strArray)-1]
	cmd := exec.Command("wget", "-O", tmpLocation, location)
	_, _ = utils.ExecuteCmd(cmd, cobraCmd)

	uploadDashboardFile(tmpLocation, cobraCmd)
}

func uploadDashboardFile(location string, cobraCmd *cobra.Command) {
	kibanaHost := viper.GetString("node.ip")
	//forcing the overwrite of already existing saved object with the same id
	cmd := exec.Command("curl", "-vX", "POST", "http://"+kibanaHost+":32003/api/kibana/dashboards/import?force=true", "-H", "Content-Type: application/json", "-H", "kbn-xsrf: true", "-d", "@"+location)
	_, _ = utils.ExecuteCmd(cmd, cobraCmd)
}
