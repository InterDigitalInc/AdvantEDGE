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
	"bufio"
	"errors"
	"os/exec"
	"strings"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

const NAMESPACE string = "NAMESPACE:"
const STATUS string = "STATUS:"
const RESOURCE string = "==>"

// Returns the status of a release
func GetReleaseStatus(name string) (*Status, error) {
	out, err := getStatus(name)
	if err != nil {
		return nil, err
	}

	status, err := parseStatus(out)
	if err != nil {
		return nil, err
	}
	return status, nil
}

func getStatus(name string) ([]byte, error) {
	var cmd = exec.Command("helm", "status", name)
	out, err := cmd.Output()
	if err != nil {
		err = errors.New("Error getting status for Release [" + name + "]")
		log.Error(err)
		return nil, err
	}
	return out, nil
}

func parseStatus(buf []byte) (*Status, error) {
	var status Status

	s := string(buf)
	scanLines := bufio.NewScanner(strings.NewReader(s))
	scanLines.Split(bufio.ScanLines)
	for i := 0; scanLines.Scan(); i++ {
		scanWords := bufio.NewScanner(strings.NewReader(scanLines.Text()))
		scanWords.Split(bufio.ScanWords)
		scanWords.Scan()
		word := scanWords.Text()

		if word == NAMESPACE {
			scanWords.Scan()
			status.Namespace = scanWords.Text()
		} else if word == STATUS {
			scanWords.Scan()
			status.State = scanWords.Text()
		} else if word == "==>" {
			var r Resource
			// Scan Type
			scanWords.Scan()
			t := strings.Split(scanWords.Text(), "/")
			r.Type = t[1]

			// Skip a line
			scanLines.Scan()

			// Scan Name
			scanLines.Scan()
			scanRes := bufio.NewScanner(strings.NewReader(scanLines.Text()))
			scanRes.Split(bufio.ScanWords)
			scanRes.Scan()
			r.Name = scanRes.Text()
			for scanRes.Scan() {
				r.Age = scanRes.Text()
			}
			status.Resources = append(status.Resources, r)
		}
	}
	return &status, nil
}
