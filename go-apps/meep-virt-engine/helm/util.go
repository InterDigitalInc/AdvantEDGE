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
	"strconv"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

type Chart struct {
	ChartName   string
	ReleaseName string
	Location    string
	ValuesFile  string
}

type Release struct {
	Name   string
	Status Status
}

type Status struct {
	State     string
	Namespace string
	Resources []Resource
}

type Resource struct {
	Name string
	Type string
	Age  string
}

func PrettyReleasesPrint(releases []Release) {
	var lines []string

	l := "#  NAME\tSTATE\t\tNAMESP.\t[RESOURCES]"
	lines = append(lines, l)
	for i, r := range releases {
		l = strconv.Itoa(i) + "- " + r.Name
		if r.Status.State != "" {
			l += "\t" + r.Status.State + "\t" + r.Status.Namespace + "\t("
			for j, res := range r.Status.Resources {
				if j != 0 {
					l += "/"
				}
				l += res.Type
			}
			l += ")"
		}
		lines = append(lines, l)
	}

	for _, l = range lines {
		log.Info(l)
	}
}
