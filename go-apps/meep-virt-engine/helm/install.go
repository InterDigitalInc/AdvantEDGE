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

package helm

import (
	"errors"
	"os/exec"
	"strings"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

func installCharts(charts []Chart) error {

	err := ensureReleases(charts)
	if err != nil {
		return err
	}

	err = install(charts)
	if err != nil {
		// Cleanup release
		cleanReleases(charts)
	}
	return err
}

func ensureReleases(charts []Chart) error {
	// ensure that releases do not already exist
	releases, _ := GetReleasesName()
	for _, c := range charts {
		for _, r := range releases {
			if c.ReleaseName == r.Name {
				err := errors.New("Release [" + c.ReleaseName + "] already exists")
				log.Error(err)
				return err
			}
		}
	}
	return nil
}

func install(charts []Chart) error {
	for _, c := range charts {
		var cmd *exec.Cmd
		if strings.Trim(c.ValuesFile, " ") == "" {
			cmd = exec.Command("helm", "install", "--name", c.ReleaseName, "--namespace", c.Namespace,
				"--set", "fullnameOverride="+c.ReleaseName, c.Location, "--replace")
		} else {
			cmd = exec.Command("helm", "install", "--name", c.ReleaseName, "--namespace", c.Namespace,
				"--set", "fullnameOverride="+c.ReleaseName, c.Location, "-f", c.ValuesFile, "--replace")
		}
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Error("Failed to install Release [" + c.ReleaseName + "] at " + c.Location)
			log.Error("Error(", err.Error(), "): ", string(out))
			return err
		}
	}
	return nil
}

func cleanReleases(charts []Chart) {
	var toClean []Chart
	var cnt int
	releases, _ := GetReleasesName()
	// ensure that releases do not exist

	for _, c := range charts {
		for _, r := range releases {
			if c.ReleaseName == r.Name {
				toClean = append(toClean, c)
				cnt++
			}
		}
	}

	if cnt > 0 {
		_ = DeleteReleases(toClean)
	}
}
