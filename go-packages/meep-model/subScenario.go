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

package model

import (
	"encoding/json"
	"net/http"

	dataModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-model"
)

// minimizeScenario - Minimizes scenario
func subScenario(scenario *dataModel.Scenario, parentLevel string, r *http.Request) ([]byte, error) {

	// Retrieve query parameters
	query := r.URL.Query()
	minimize := query.Get("minimize")
	children := query.Get("children")
	domainName := query.Get("domain")
	domainType := query.Get("domainType")
	zoneName := query.Get("zone")
	nlName := query.Get("nl")
	nlType := query.Get("nlType")
	plName := query.Get("pl")
	plType := query.Get("plType")
	procName := query.Get("proc")
	procType := query.Get("procType")
	var domains []dataModel.Domain
	var zones []dataModel.Zone
	var nls []dataModel.NetworkLocation
	var pls []dataModel.PhysicalLocation
	var procs []dataModel.Process
	var foundDomain, validDomainChildren bool
	var foundZone, validZoneChildren bool
	var foundNl, validNlChildren bool
	var foundPl, validPlChildren bool
	//var foundProc, validProcChildren bool

	if minimize == "true" {

		err := minimizeScenario(scenario)
		if err != nil {
			return nil, err
		}
	}

	if scenario != nil {
		if scenario.Deployment != nil {
			deployment := scenario.Deployment

			// Domains
			for iDomain := range deployment.Domains {
				domain := &deployment.Domains[iDomain]
				foundDomain = true

				if domainName != "" {
					if domainName != domain.Name {
						foundDomain = false
						continue
					}
				}
				if domainType != "" {
					if domainType != domain.Type_ {
						foundDomain = false
						continue
					}
				}

				// Zones
				if len(domain.Zones) > 0 || zoneName != "" || nlName != "" || plName != "" || procName != "" || nlType != "" || plType != "" || procType != "" {
					validDomainChildren = false
				} else {
					validDomainChildren = true
				}

				for iZone := range domain.Zones {
					zone := &domain.Zones[iZone]
					foundZone = true
					if zoneName != "" {
						if zoneName != zone.Name {
							foundZone = false
							continue
						}
					}
					// Network Locations
					if len(zone.NetworkLocations) > 0 || nlName != "" || plName != "" || procName != "" || nlType != "" || plType != "" || procType != "" {
						validZoneChildren = false
					} else {
						validZoneChildren = true
					}

					for iNL := range zone.NetworkLocations {
						nl := &zone.NetworkLocations[iNL]
						foundNl = true

						if nlName != "" {
							if nlName != nl.Name {
								foundNl = false
								continue
							}
						}
						if nlType != "" {
							if nlType != nl.Type_ {
								foundNl = false
								continue
							}
						}

						// Physical Locations
						if len(nl.PhysicalLocations) > 0 || plName != "" || procName != "" || plType != "" || procType != "" {
							validNlChildren = false
						} else {
							validNlChildren = true
						}

						for iPL := range nl.PhysicalLocations {
							pl := &nl.PhysicalLocations[iPL]
							foundPl = true

							if plName != "" {
								if plName != pl.Name {
									foundPl = false
									continue
								}
							}
							if plType != "" {
								if plType != pl.Type_ {
									foundPl = false
									continue
								}
							}
							// Processes
							if len(pl.Processes) > 0 || procName != "" || procType != "" {
								validPlChildren = false
							} else {
								validPlChildren = true
							}

							for iProc := range pl.Processes {
								proc := &pl.Processes[iProc]
								if procName != "" {
									if procName != proc.Name {
										continue
									}
								}
								if procType != "" {
									if procType != proc.Type_ {
										continue
									}
								}

								if parentLevel == "proc" {
									procs = append(procs, *proc)
								}
								validPlChildren = true
							}

							if foundPl && validPlChildren {
								if parentLevel == "pl" {
									if children == "false" {
										pl.Processes = nil
									}
									pls = append(pls, *pl)
								}
								validNlChildren = true
							}
						}

						if foundNl && validNlChildren {
							if parentLevel == "nl" {
								if children == "false" {
									nl.PhysicalLocations = nil
								}
								nls = append(nls, *nl)
							}
							validZoneChildren = true
						}
					}
					if foundZone && validZoneChildren {
						if parentLevel == "zone" {
							if children == "false" {
								zone.NetworkLocations = nil
							}
							zones = append(zones, *zone)
						}
						validDomainChildren = true
					}
				}
				if foundDomain && validDomainChildren && parentLevel == "domain" {
					if children == "false" {
						domain.Zones = nil
					}
					domains = append(domains, *domain)
				}
			}
		}
	}
	var result []byte
	var err error
	switch parentLevel {
	case "domain":
		result, err = json.Marshal(domains)
	case "zone":
		result, err = json.Marshal(zones)
	case "nl":
		result, err = json.Marshal(nls)
	case "pl":
		result, err = json.Marshal(pls)
	case "proc":
		result, err = json.Marshal(procs)

	default:
	}

	if err != nil {
		return nil, err
	}
	return result, nil
}
