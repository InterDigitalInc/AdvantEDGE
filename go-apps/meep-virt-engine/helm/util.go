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
	"strconv"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

type Chart struct {
	Name        string
	ReleaseName string
	Namespace   string
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
