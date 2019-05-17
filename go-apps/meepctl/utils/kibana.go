/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */

package utils

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/roymx/viper"
	"github.com/spf13/cobra"
)

func DeployKibanaDashboards(cobraCmd *cobra.Command) {

	gitDir := viper.GetString("meep.gitdir") + "/"
	workdir := viper.GetString("meep.workdir")

	start := time.Now()

	//make sure kibana is up and ready to receive messages
	//this only happens when kibana is able to connect to elastic search
	//so elastic search must be up and running as well as kibana for the rest api server to be up
	isKibanaUp := false
	kibanaFailedAttempts := 0
	for !isKibanaUp {
		isKibanaUp = isKibanaReady(cobraCmd)
		if !isKibanaUp {
			kibanaFailedAttempts++
			if kibanaFailedAttempts > 3 {
				elapsed := time.Since(start)
				r := FormatResult("Failure during deployment kibana dashboards", elapsed, cobraCmd)
				fmt.Println(r)
				return
			}
		}
	}

	cmd := exec.Command("cp", gitDir+"dashboards/dashboards.conf", workdir+"/tmp/dashboards.conf")
	_, _ = ExecuteCmd(cmd, cobraCmd)
	//search and replace in yaml file
	tmpStr := strings.Replace(gitDir, "/", "\\/", -1)
	str := "s/<GITDIR>/" + tmpStr + "/g"
	cmd = exec.Command("sed", "-i", str, workdir+"/tmp/dashboards.conf")
	_, _ = ExecuteCmd(cmd, cobraCmd)

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
	r := FormatResult("Deployed kibana dashboards", elapsed, cobraCmd)
	fmt.Println(r)

}

//sending a DUMMY value just to see if the service is up (conditions to be up are:
//- all elastic search(ES) pods are up
//- kibana pod is up
//- kibana connected successfully to ES
func isKibanaReady(cobraCmd *cobra.Command) bool {
	isReady := false
	err := uploadDefaultIndex("DUMMY", cobraCmd)
	if err == nil {
		isReady = true
	}
	return isReady
}

func uploadDefaultIndex(indexId string, cobraCmd *cobra.Command) error {
	kibanaHost, _ := os.Hostname()
	cmd := exec.Command("curl", "-vX", "POST", "http://"+kibanaHost+":32003/api/kibana/settings/defaultIndex", "-H", "Content-Type: application/json", "-H", "kbn-xsrf: true", "-d", "{\"value\": \""+indexId+"\"}")
	out, err := cmd.CombinedOutput()
	if err != nil {
		err = errors.New("Error sending a curl command")
	} else {
		str := string(out)
		isServiceUnavailable := strings.Contains(str, "Service Unavailable")

		if isServiceUnavailable {
			err = errors.New("Error: Service Unavailable")
		}
	}

	return err
}

func uploadDashboardHttp(location string, cobraCmd *cobra.Command) {
	//no support yet for url in kibana 6.4.2... but we can get the file to /tmp and then download the file
	strArray := strings.Split(location, "/")
	tmpLocation := "/tmp/" + strArray[len(strArray)-1]
	cmd := exec.Command("wget", "-O", tmpLocation, location)
	_, _ = ExecuteCmd(cmd, cobraCmd)

	uploadDashboardFile(tmpLocation, cobraCmd)
}

func uploadDashboardFile(location string, cobraCmd *cobra.Command) {
	kibanaHost, _ := os.Hostname()
	//forcing the overwrite of already existing saved object with the same id
	cmd := exec.Command("curl", "-vX", "POST", "http://"+kibanaHost+":32003/api/kibana/dashboards/import?force=true", "-H", "Content-Type: application/json", "-H", "kbn-xsrf: true", "-d", "@"+location)
	_, _ = ExecuteCmd(cmd, cobraCmd)
}
