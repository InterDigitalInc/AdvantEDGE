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

package helm

import (
	"bufio"
	"errors"
	"os/exec"
	"strings"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

func getReleasesName(sandboxName string) ([]Release, error) {
	out, err := getList(sandboxName)
	if err != nil {
		return nil, err
	}

	release, err := parseList(out, true, sandboxName)
	if err != nil {
		return nil, err
	}
	return release, nil
}

func getReleases(sandboxName string) ([]Release, error) {
	out, err := getList(sandboxName)
	if err != nil {
		return nil, err
	}

	release, err := parseList(out, false, sandboxName)
	if err != nil {
		return nil, err
	}
	return release, nil
}

func getList(sandboxName string) ([]byte, error) {
	var cmd = exec.Command("helm", "ls", "-n", sandboxName)
	out, err := cmd.Output()
	if err != nil {
		err = errors.New("Unable to list Releases")
		log.Error(err)
		return nil, err
	}
	return out, nil
}

func parseList(buf []byte, nameOnly bool, sandboxName string) ([]Release, error) {
	/* Example of what needs to be parsed
	NAME    REVISION        UPDATED                         STATUS          CHART                   NAMESPACE
	osvc1   1               Tue Jun 12 13:02:55 2018        DEPLOYED        orientation-svc-0.1.0   default
	osvc2   1               Tue Jun 12 13:04:54 2018        DEPLOYED        orientation-svc-0.1.0   default
	uss     1               Wed Jun 13 19:42:29 2018        DEPLOYED        uss-0.1.0               default
	uss-db  1               Wed Jun 13 19:41:34 2018        DEPLOYED        mysql-0.5.0             default
	*/
	var releases []Release

	s := string(buf)
	scanLines := bufio.NewScanner(strings.NewReader(s))
	scanLines.Split(bufio.ScanLines)
	for i := 0; scanLines.Scan(); i++ {
		if i == 0 {
			continue
		}
		scanWords := bufio.NewScanner(strings.NewReader(scanLines.Text()))
		scanWords.Split(bufio.ScanWords)
		scanWords.Scan()
		var r Release
		// Name
		r.Name = scanWords.Text()
		if !nameOnly {
			// Status
			sp, err := GetReleaseStatus(r.Name, sandboxName)
			r.Status = *sp
			if err != nil {
				log.Error(err)
				continue
			}
		}
		releases = append(releases, r)
	}

	return releases, nil
}
