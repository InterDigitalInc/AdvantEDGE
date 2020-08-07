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
	dataModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-model"
)

// minimizeScenario - Minimizes scenario
func minimizeScenario(scenario *dataModel.Scenario) error {
	if scenario != nil {
		if scenario.Deployment != nil {
			deployment := scenario.Deployment

			// Domains
			for iDomain := range deployment.Domains {
				domain := &deployment.Domains[iDomain]

				// Zones
				for iZone := range domain.Zones {
					zone := &domain.Zones[iZone]

					// Network Locations
					for iNL := range zone.NetworkLocations {
						nl := &zone.NetworkLocations[iNL]

						// Remove geodata
						nl.GeoData = nil

						// Physical Locations
						for iPL := range nl.PhysicalLocations {
							pl := &nl.PhysicalLocations[iPL]

							// Remove geodata
							pl.GeoData = nil

							// // Processes
							// for iProc := range pl.Processes {
							// 	proc := &pl.Processes[iProc]
							// }
						}
					}
				}
			}
		}
	}
	return nil
}
