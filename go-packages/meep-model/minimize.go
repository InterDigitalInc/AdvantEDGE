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
func minimizeScenario(scenario *dataModel.Scenario) {
	if scenario != nil {
		minimizeDeployment(scenario.Deployment)
	}
}

// minimizeDeployment - Minimizes deployment
func minimizeDeployment(deployment *dataModel.Deployment) {
	if deployment != nil {
		for i := range deployment.Domains {
			domain := &deployment.Domains[i]
			minimizeDomain(domain)
		}
	}
}

// minimizeDomain - Minimizes domain
func minimizeDomain(domain *dataModel.Domain) {
	for i := range domain.Zones {
		zone := &domain.Zones[i]
		minimizeZone(zone)
	}
}

// minimizeZone - Minimizes zone
func minimizeZone(zone *dataModel.Zone) {
	for i := range zone.NetworkLocations {
		nl := &zone.NetworkLocations[i]
		minimizeNetLoc(nl)
	}
}

// minimizeNetLoc - Minimizes network location
func minimizeNetLoc(nl *dataModel.NetworkLocation) {
	// Remove geodata
	nl.GeoData = nil

	for i := range nl.PhysicalLocations {
		pl := &nl.PhysicalLocations[i]
		minimizePhyLoc(pl)
	}
}

// minimizePlyLoc - Minimizes physical location
func minimizePhyLoc(pl *dataModel.PhysicalLocation) {
	// Remove geodata
	pl.GeoData = nil
}
