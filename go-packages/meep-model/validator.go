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
	"errors"

	dataModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-model"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"

	"github.com/blang/semver"
)

// Validator status types
const (
	ValidatorStatusValid   = "SCENARIO-VALID"
	ValidatorStatusUpdated = "SCENARIO-UPDATED"
	ValidatorStatusError   = "SCENARIO-ERROR"
)

var DefaultVersion = semver.Version{Major: 1, Minor: 0, Patch: 0}
var ValidatorVersion = semver.Version{Major: 1, Minor: 4, Patch: 0}

// ValidateScenario - Verify if json scenario is valid & supported. Upgrade scenario if possible & necessary.
func ValidateScenario(jsonScenario []byte) (validJsonScenario []byte, status string, err error) {
	var scenarioVersion semver.Version

	// Unmarshal scenario
	scenario := new(dataModel.Scenario)
	err = json.Unmarshal(jsonScenario, scenario)
	if err != nil {
		log.Error(err.Error())
		return nil, ValidatorStatusError, err
	}

	// Retrieve scenario version
	if scenario.Version == "" {
		scenarioVersion = DefaultVersion
	} else {
		scenarioVersion, err = semver.Make(scenario.Version)
		if err != nil {
			log.Error(err.Error())
			return nil, ValidatorStatusError, err
		}
	}

	// Verify that scenario is compatible
	if scenarioVersion.Major != ValidatorVersion.Major || scenarioVersion.GT(ValidatorVersion) {
		err = errors.New("Scenario version " + scenario.Version + " incompatible with validator version " + ValidatorVersion.String())
		return nil, ValidatorStatusError, err
	}

	// Upgrade scenario if necessary
	if scenarioVersion.EQ(ValidatorVersion) {
		return jsonScenario, ValidatorStatusValid, nil
	} else {
		// Set updated version
		scenario.Version = ValidatorVersion.String()

		// Migrate zone information
		if scenario.Deployment != nil {
			for iDomain := range scenario.Deployment.Domains {
				domain := &scenario.Deployment.Domains[iDomain]
				for iZone := range domain.Zones {
					zone := &domain.Zones[iZone]

					// Create new Network Characteristic field and migrate values from EdgeFog
					if zone.NetChar == nil {
						zone.NetChar = new(dataModel.NetworkCharacteristics)
						zone.NetChar.Latency = zone.EdgeFogLatency
						zone.NetChar.LatencyVariation = zone.EdgeFogLatencyVariation
						zone.NetChar.PacketLoss = zone.EdgeFogPacketLoss
						zone.NetChar.Throughput = zone.EdgeFogThroughput
					}

					// Reset deprecated values to omit them
					zone.InterEdgeLatency = 0
					zone.InterEdgeLatencyVariation = 0
					zone.InterEdgePacketLoss = 0
					zone.InterEdgeThroughput = 0
					zone.InterFogLatency = 0
					zone.InterFogLatencyVariation = 0
					zone.InterFogPacketLoss = 0
					zone.InterFogThroughput = 0
					zone.EdgeFogLatency = 0
					zone.EdgeFogLatencyVariation = 0
					zone.EdgeFogPacketLoss = 0
					zone.EdgeFogThroughput = 0
				}
			}
		}
	}

	// Marshal updated scenario
	validJsonScenario, err = json.Marshal(scenario)
	if err != nil {
		return nil, ValidatorStatusError, err
	}
	return validJsonScenario, ValidatorStatusUpdated, err
}

// Validate the provided PL
func validatePL(pl *dataModel.PhysicalLocation) error {

	if pl.Id == "" {
		pl.Id = pl.Name
	}
	if pl.Name == "" {
		return errors.New("Invalid Name")
	}
	if pl.Type_ != NodeTypeUE {
		return errors.New("Unsupported PL Type: " + pl.Type_)
	}
	if pl.LinkThroughput == 0 {
		pl.LinkThroughput = 1000
	}

	return nil
}
