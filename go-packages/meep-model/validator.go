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

// Current validator version
var ValidatorVersion = semver.Version{Major: 1, Minor: 6, Patch: 0}

// Versions requiring scenario update
var Version130 = semver.Version{Major: 1, Minor: 3, Patch: 0}
var Version140 = semver.Version{Major: 1, Minor: 4, Patch: 0}
var Version150 = semver.Version{Major: 1, Minor: 5, Patch: 0}
var Version160 = semver.Version{Major: 1, Minor: 6, Patch: 0}

// Default latency distribution
const DEFAULT_LATENCY_DISTRIBUTION = "Normal"

// setNetChar - Creates a new netchar object if non-existent and migrate values from deprecated fields
func createNetChar(lat int32, latVar int32, dist string, tputDl int32, tputUl int32, loss float64) *dataModel.NetworkCharacteristics {
	nc := new(dataModel.NetworkCharacteristics)
	nc.Latency = lat
	nc.LatencyVariation = latVar
	nc.LatencyDistribution = dist
	nc.PacketLoss = loss
	nc.ThroughputDl = tputDl
	nc.ThroughputUl = tputUl
	return nc
}

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
	// If no version found, assume current validator version
	if scenario.Version == "" {
		scenarioVersion = ValidatorVersion
	} else {
		scenarioVersion, err = semver.Make(scenario.Version)
		if err != nil {
			log.Error(err.Error())
			return nil, ValidatorStatusError, err
		}

		// Verify that scenario is compatible
		if scenarioVersion.Major != ValidatorVersion.Major ||
			scenarioVersion.GT(ValidatorVersion) ||
			scenarioVersion.LT(Version130) {
			err = errors.New("Scenario version " + scenario.Version + " incompatible with validator version " + ValidatorVersion.String())
			return nil, ValidatorStatusError, err
		}

		// Skip validation if already current version
		if scenarioVersion.EQ(ValidatorVersion) {
			return jsonScenario, ValidatorStatusValid, nil
		}
	}

	// Run upgrade functions starting from oldest applicable patch to newest

	// UPGRADE TO 1.4.0
	if scenarioVersion.LT(Version140) {
		upgradeScenarioTo140(scenario)
		scenarioVersion = Version140
	}
	// UPGRADE TO 1.5.0
	if scenarioVersion.LT(Version150) {
		upgradeScenarioTo150(scenario)
		scenarioVersion = Version150
	}
	// UPGRADE TO 1.6.0
	if scenarioVersion.LT(Version160) {
		upgradeScenarioTo160(scenario)
		scenarioVersion = Version160
	}

	// Set current scenario version
	scenario.Version = ValidatorVersion.String()

	// Marshal updated scenario
	validJsonScenario, err = json.Marshal(scenario)
	if err != nil {
		return nil, ValidatorStatusError, err
	}
	return validJsonScenario, ValidatorStatusUpdated, err
}

func upgradeScenarioTo140(scenario *dataModel.Scenario) {
	// Set updated version
	scenario.Version = Version140.String()

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

func upgradeScenarioTo150(scenario *dataModel.Scenario) {
	// Set updated version
	scenario.Version = Version150.String()

	// Migrate netchar information
	if scenario.Deployment != nil {
		deploy := scenario.Deployment

		// Create new Network Characteristic field and migrate values, if necessary
		if deploy.NetChar == nil {
			deploy.NetChar = createNetChar(
				deploy.InterDomainLatency,
				deploy.InterDomainLatencyVariation,
				DEFAULT_LATENCY_DISTRIBUTION,
				deploy.InterDomainThroughput,
				deploy.InterDomainThroughput,
				deploy.InterDomainPacketLoss)
		}

		// Reset deprecated values to omit them
		deploy.InterDomainLatency = 0
		deploy.InterDomainLatencyVariation = 0
		deploy.InterDomainPacketLoss = 0
		deploy.InterDomainThroughput = 0

		for iDomain := range scenario.Deployment.Domains {
			domain := &scenario.Deployment.Domains[iDomain]

			// Create new Network Characteristic field and migrate values, if necessary
			if domain.NetChar == nil {
				domain.NetChar = createNetChar(
					domain.InterZoneLatency,
					domain.InterZoneLatencyVariation,
					"",
					domain.InterZoneThroughput,
					domain.InterZoneThroughput,
					domain.InterZonePacketLoss)
			}

			// Reset deprecated values to omit them
			domain.InterZoneLatency = 0
			domain.InterZoneLatencyVariation = 0
			domain.InterZonePacketLoss = 0
			domain.InterZoneThroughput = 0

			for iZone := range domain.Zones {
				zone := &domain.Zones[iZone]

				// Migrate throughput values, if necessary
				if zone.NetChar.ThroughputDl == 0 {
					zone.NetChar.ThroughputDl = zone.NetChar.Throughput
					zone.NetChar.ThroughputUl = zone.NetChar.Throughput
				}

				// Reset deprecated values to omit
				zone.NetChar.Throughput = 0

				for iNl := range zone.NetworkLocations {
					nl := &zone.NetworkLocations[iNl]

					// Create new Network Characteristic field and migrate values, if necessary
					if nl.NetChar == nil {
						nl.NetChar = createNetChar(
							nl.TerminalLinkLatency,
							nl.TerminalLinkLatencyVariation,
							"",
							nl.TerminalLinkThroughput,
							nl.TerminalLinkThroughput,
							nl.TerminalLinkPacketLoss)
					}

					// Reset deprecated values to omit them
					nl.TerminalLinkLatency = 0
					nl.TerminalLinkLatencyVariation = 0
					nl.TerminalLinkPacketLoss = 0
					nl.TerminalLinkThroughput = 0

					// Physical Locations
					for iPl := range nl.PhysicalLocations {
						pl := &nl.PhysicalLocations[iPl]

						// Create new Network Characteristic field and migrate values, if necessary
						if pl.NetChar == nil {
							pl.NetChar = createNetChar(
								pl.LinkLatency,
								pl.LinkLatencyVariation,
								"",
								pl.LinkThroughput,
								pl.LinkThroughput,
								pl.LinkPacketLoss)
						}

						// Reset deprecated values to omit them
						pl.LinkLatency = 0
						pl.LinkLatencyVariation = 0
						pl.LinkPacketLoss = 0
						pl.LinkThroughput = 0

						for iProc := range pl.Processes {
							proc := &pl.Processes[iProc]

							// Create new Network Characteristic field and migrate values, if necessary
							if proc.NetChar == nil {
								proc.NetChar = createNetChar(
									proc.AppLatency,
									proc.AppLatencyVariation,
									"",
									proc.AppThroughput,
									proc.AppThroughput,
									proc.AppPacketLoss)
							}

							// Reset deprecated values to omit them
							proc.AppLatency = 0
							proc.AppLatencyVariation = 0
							proc.AppPacketLoss = 0
							proc.AppThroughput = 0
						}
					}
				}
			}
		}
	}
}

func upgradeScenarioTo160(scenario *dataModel.Scenario) {
	//changes in 160 vs 150
	//rename POA-CELLULAR to POA-4G

	// Set updated version
	scenario.Version = Version160.String()

	// Migrate netchar information
	if scenario.Deployment != nil {
		for iDomain := range scenario.Deployment.Domains {
			domain := &scenario.Deployment.Domains[iDomain]
			for iZone := range domain.Zones {
				zone := &domain.Zones[iZone]
				for iNl := range zone.NetworkLocations {
					nl := &zone.NetworkLocations[iNl]
					if nl.Type_ == NodeTypePoaCellularDeprecated /*"POA-CELLULAR"*/ {
						nl.Type_ = NodeTypePoa4G
						if nl.CellularPoaConfig != nil {
							if nl.Poa4GConfig == nil {
								nl.Poa4GConfig = new(dataModel.Poa4GConfig)
							}
							nl.Poa4GConfig.CellId = nl.CellularPoaConfig.CellId
						}
					}
				}
			}
		}
	}
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
	if pl.NetChar != nil {
		if pl.NetChar.ThroughputDl == 0 {
			pl.NetChar.ThroughputDl = 1000
		}
		if pl.NetChar.ThroughputUl == 0 {
			pl.NetChar.ThroughputUl = 1000
		}
	}
	return nil
}
