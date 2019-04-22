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
	"encoding/json"
	"strconv"

	log "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-virt-engine/log"
)

const (
	ChartMeepType      = "MEEP-TYPE"
	ChartUserChartType = "USERCHART-TYPE"
)

type Chart struct {
	Type            string
	ChartName       string
	ReleaseName     string
	Location        string
	AlternateValues string
	Parameters      string
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

func prettyJsonPrint(v interface{}) {
	j, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		log.Error(err)
		return
	}
	log.Info(string(j))
}

func PrettyReleasesPrint(releases []Release) {
	var lines []string
	var l string

	l = "#  NAME\tSTATE\t\tNAMESP.\t[RESOURCES]"
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
