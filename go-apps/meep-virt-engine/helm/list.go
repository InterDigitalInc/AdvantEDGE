/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */
package helm

import (
	"bufio"
	"errors"
	"os/exec"
	"strings"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

func getReleasesName() ([]Release, error) {
	out, err := getList()
	if err != nil {
		return nil, err
	}

	release, err := parseList(out, true)
	if err != nil {
		return nil, err
	}
	return release, nil
}

func getReleases() ([]Release, error) {
	out, err := getList()
	if err != nil {
		return nil, err
	}

	release, err := parseList(out, false)
	if err != nil {
		return nil, err
	}
	return release, nil
}

func getList() ([]byte, error) {
	var cmd = exec.Command("helm", "ls")
	out, err := cmd.Output()
	if err != nil {
		err = errors.New("Unable to list Releases")
		log.Error(err)
		return nil, err
	}
	return out, nil
}

func parseList(buf []byte, nameOnly bool) ([]Release, error) {
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
			sp, err := GetReleaseStatus(r.Name)
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
