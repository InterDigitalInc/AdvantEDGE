/*
 * Copyright (c) 2022  The AdvantEDGE Authors
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

package meepdaimgr

import (
	"errors"
	"os/exec"
	"strings"
	"time"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"

	uuid "github.com/google/uuid"
	"github.com/spf13/cobra"
)

var cbCmd = &cobra.Command{
	Use:       "Testing exec",
	Short:     "Testing exec",
	Long:      "Long description",
	Example:   "Usage example",
	Args:      cobra.OnlyValidArgs,
	ValidArgs: nil,
	Run:       nil,
}

// Create and start a container
func DockerRun(image string, port string, timeout time.Duration) (dockerId string, out string, err error) {
	// log.Debug(">>> DockerRun: port: ", port)
	// log.Debug(">>> DockerRun: image: ", image)
	// log.Debug(">>> DockerRun: timeout: ", timeout)

	// Sanity checks
	if image == "" {
		err = errors.New("Wrong parameters")
		return "", "", err
	}
	cmd := exec.Command("docker", "images", image)
	out, err = ExecuteCmd(cmd, cbCmd)
	if err != nil {
		log.Error(err.Error())
		return "", out, err
	}
	log.Debug("DockerRun: out: ", out)
	if len(strings.Split(out, "\n")) < 2 {
		err = errors.New("Image does not exist in repository")
		log.Error(err.Error())
		return "", out, err

	}

	dockerId = uuid.New().String()[0:32] // length limited to 32 chars, see MEC-016 Clause 6.2.3 Type: AppContext
	if port == "" {
		cmd = exec.Command("docker", "run", "--name", dockerId, "--rm", "-d", "-t", image)
	} else {
		cmd = exec.Command("docker", "run", "--name", dockerId, "--rm", "--expose="+port, "-d", "-t", image)
	}
	out, err = ExecuteCmd(cmd, cbCmd)
	if err != nil {
		log.Error(err.Error())
		return "", out, err
	}
	log.Debug("DockerRun: dockerId: ", dockerId)
	log.Debug("DockerRun: out: ", out)

	time.Sleep(timeout) // Wait for container starting

	return dockerId, out, nil
}

// Create and start a container
func DockerTerminate(dockerId string, timeout time.Duration) (out string, err error) {
	//log.Debug(">>> DockerTerminate: ", dockerId)

	// Sanity checks
	if dockerId == "" {
		err = errors.New("Wromg parameters")
		return "", err
	}

	cmd := exec.Command("docker", "kill", dockerId)
	out, err = ExecuteCmd(cmd, cbCmd)
	if err != nil {
		log.Error(err.Error())
		return out, err
	}
	log.Debug("DockerRun: dockerId: ", dockerId)
	log.Debug("DockerRun: out: ", out)

	time.Sleep(timeout) // Wait for container starting

	return out, nil
}
