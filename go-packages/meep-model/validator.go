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
var ValidatorVersion = semver.Version{Major: 1, Minor: 5, Patch: 0}

// Default latencies per physical location type
const DEFAULT_LATENCY_INTER_DOMAIN = 50
const DEFAULT_LATENCY_JITTER_INTER_DOMAIN = 10
const DEFAULT_LATENCY_DISTRIBUTION_INTER_DOMAIN = "Normal"
const DEFAULT_THROUGHPUT_DL_INTER_DOMAIN = 1000
const DEFAULT_THROUGHPUT_UL_INTER_DOMAIN = 1000
const DEFAULT_PACKET_LOSS_INTER_DOMAIN = 0
const DEFAULT_LATENCY_INTER_ZONE = 6
const DEFAULT_LATENCY_JITTER_INTER_ZONE = 2
const DEFAULT_THROUGHPUT_DL_INTER_ZONE = 1000
const DEFAULT_THROUGHPUT_UL_INTER_ZONE = 1000
const DEFAULT_PACKET_LOSS_INTER_ZONE = 0
const DEFAULT_LATENCY_INTRA_ZONE = 5
const DEFAULT_LATENCY_JITTER_INTRA_ZONE = 1
const DEFAULT_THROUGHPUT_DL_INTRA_ZONE = 1000
const DEFAULT_THROUGHPUT_UL_INTRA_ZONE = 1000
const DEFAULT_PACKET_LOSS_INTRA_ZONE = 0
const DEFAULT_LATENCY_TERMINAL_LINK = 1
const DEFAULT_LATENCY_JITTER_TERMINAL_LINK = 1
const DEFAULT_THROUGHPUT_DL_TERMINAL_LINK = 1000
const DEFAULT_THROUGHPUT_UL_TERMINAL_LINK = 1000
const DEFAULT_PACKET_LOSS_TERMINAL_LINK = 0
const DEFAULT_LATENCY_LINK = 0
const DEFAULT_LATENCY_JITTER_LINK = 0
const DEFAULT_THROUGHPUT_DL_LINK = 1000
const DEFAULT_THROUGHPUT_UL_LINK = 1000
const DEFAULT_PACKET_LOSS_LINK = 0
const DEFAULT_LATENCY_APP = 0
const DEFAULT_LATENCY_JITTER_APP = 0
const DEFAULT_THROUGHPUT_DL_APP = 1000
const DEFAULT_THROUGHPUT_UL_APP = 1000
const DEFAULT_PACKET_LOSS_APP = 0
const DEFAULT_LATENCY_DC = 0

// setNetChar - Creates a new netchar object if non-existent and migrate values from deprecated fields
func createNetChar(latency int32, latencyVariation int32, distribution string, throughputDl int32, throughputUl int32, packetLoss float64) *dataModel.NetworkCharacteristics {

	nc := new(dataModel.NetworkCharacteristics)
	nc.Latency = latency
	nc.LatencyVariation = latencyVariation
	nc.LatencyDistribution = distribution
	nc.PacketLoss = packetLoss
	nc.ThroughputDl = throughputDl
	nc.ThroughputUl = throughputUl
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
		previousVersion := scenario.Version
		scenario.Version = ValidatorVersion.String()

		//NetChar validator variables
		var latency int32
		var latencyVariation int32
		var latencyDistribution string
		var throughputDl int32
		var throughputUl int32
		var packetLoss float64

		// Migrate netchar information
		if scenario.Deployment != nil {
			deploy := scenario.Deployment
			latency = deploy.InterDomainLatency
			latencyVariation = deploy.InterDomainLatencyVariation
			latencyDistribution = DEFAULT_LATENCY_DISTRIBUTION_INTER_DOMAIN
			throughputDl = deploy.InterDomainThroughput
			throughputUl = deploy.InterDomainThroughput
			packetLoss = deploy.InterDomainPacketLoss

			nc := deploy.NetChar
			if nc != nil {
				//netchar got already created, if values are default, replace them with what was in the scenario, otherwise, leave them as is
				if nc.Latency != DEFAULT_LATENCY_INTER_DOMAIN {
					latency = nc.Latency
				}
				if nc.LatencyVariation != DEFAULT_LATENCY_JITTER_INTER_DOMAIN {
					latencyVariation = nc.LatencyVariation
				}
				if nc.LatencyDistribution != DEFAULT_LATENCY_DISTRIBUTION_INTER_DOMAIN {
					latencyDistribution = nc.LatencyDistribution
				}
				if nc.ThroughputDl != DEFAULT_THROUGHPUT_DL_INTER_DOMAIN {
					throughputDl = nc.ThroughputDl
				}
				if nc.ThroughputUl != DEFAULT_THROUGHPUT_UL_INTER_DOMAIN {
					throughputUl = nc.ThroughputUl
				}
				if nc.PacketLoss != DEFAULT_PACKET_LOSS_INTER_DOMAIN {
					packetLoss = nc.PacketLoss
				}
			}
			deploy.NetChar = createNetChar(latency, latencyVariation, latencyDistribution, throughputDl, throughputUl, packetLoss)

			// Reset deprecated values to omit them
			deploy.InterDomainLatency = 0
			deploy.InterDomainLatencyVariation = 0
			deploy.InterDomainPacketLoss = 0
			deploy.InterDomainThroughput = 0

			for iDomain := range scenario.Deployment.Domains {
				domain := &scenario.Deployment.Domains[iDomain]
				latency = domain.InterZoneLatency
				latencyVariation = domain.InterZoneLatencyVariation
				throughputDl = domain.InterZoneThroughput
				throughputUl = domain.InterZoneThroughput
				packetLoss = domain.InterZonePacketLoss
				nc := domain.NetChar
				if nc != nil {
					//netchar got already created, if values are default, replace them with what was in the scenario, otherwise, leave them as is
					if nc.Latency != DEFAULT_LATENCY_INTER_ZONE {
						latency = nc.Latency
					}
					if nc.LatencyVariation != DEFAULT_LATENCY_JITTER_INTER_ZONE {
						latencyVariation = nc.LatencyVariation
					}
					if nc.ThroughputDl != DEFAULT_THROUGHPUT_DL_INTER_ZONE {
						throughputDl = nc.ThroughputDl
					}
					if nc.ThroughputUl != DEFAULT_THROUGHPUT_UL_INTER_ZONE {
						throughputUl = nc.ThroughputUl
					}
					if nc.PacketLoss != DEFAULT_PACKET_LOSS_INTER_ZONE {
						packetLoss = nc.PacketLoss
					}
				}
				domain.NetChar = createNetChar(latency, latencyVariation, "", throughputDl, throughputUl, packetLoss)

				// Reset deprecated values to omit them
				domain.InterZoneLatency = 0
				domain.InterZoneLatencyVariation = 0
				domain.InterZonePacketLoss = 0
				domain.InterZoneThroughput = 0

				for iZone := range domain.Zones {
					zone := &domain.Zones[iZone]
					latency = zone.EdgeFogLatency
					latencyVariation = zone.EdgeFogLatencyVariation
					throughputDl = zone.EdgeFogThroughput
					throughputUl = zone.EdgeFogThroughput
					packetLoss = zone.EdgeFogPacketLoss
					if zone.NetChar != nil {
						nc := zone.NetChar
						//migration from 1.3
						if previousVersion == "1.3" {
							latency = zone.EdgeFogLatency
							latencyVariation = zone.EdgeFogLatencyVariation
							throughputDl = zone.EdgeFogThroughput
							throughputUl = zone.EdgeFogThroughput
							packetLoss = zone.EdgeFogPacketLoss
						} else {
							latency = nc.Latency
							latencyVariation = nc.LatencyVariation
							throughputDl = nc.Throughput
							throughputUl = nc.Throughput
							packetLoss = nc.PacketLoss
						}
						//netchar got already created, if values are default, replace them with what was in the scenario, otherwise, leave them as is
						if nc.Latency != DEFAULT_LATENCY_INTER_ZONE {
							latency = nc.Latency
						}
						if nc.LatencyVariation != DEFAULT_LATENCY_JITTER_INTER_ZONE {
							latencyVariation = nc.LatencyVariation
						}
						if nc.ThroughputDl != DEFAULT_THROUGHPUT_DL_INTER_ZONE {
							throughputDl = nc.ThroughputDl
						}
						if nc.ThroughputUl != DEFAULT_THROUGHPUT_UL_INTER_ZONE {
							throughputUl = nc.ThroughputUl
						}
						if nc.PacketLoss != DEFAULT_PACKET_LOSS_INTER_ZONE {
							packetLoss = nc.PacketLoss
						}
					}
					zone.NetChar = createNetChar(latency, latencyVariation, "", throughputDl, throughputUl, packetLoss)
					// Reset deprecated values to omit them from v1.4
					zone.NetChar.Throughput = 0
					// Reset deprecated values to omit them from v1.3
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
					for iNl := range zone.NetworkLocations {
						nl := &zone.NetworkLocations[iNl]
						latency = nl.TerminalLinkLatency
						latencyVariation = nl.TerminalLinkLatencyVariation
						throughputDl = nl.TerminalLinkThroughput
						throughputUl = nl.TerminalLinkThroughput
						packetLoss = nl.TerminalLinkPacketLoss
						nc := nl.NetChar
						if nc != nil {
							//netchar got already created, if values are default, replace them with what was in the scenario, otherwise, leave them as is
							if nc.Latency != DEFAULT_LATENCY_TERMINAL_LINK {
								latency = nc.Latency
							}
							if nc.LatencyVariation != DEFAULT_LATENCY_JITTER_TERMINAL_LINK {
								latencyVariation = nc.LatencyVariation
							}
							if nc.ThroughputDl != DEFAULT_THROUGHPUT_DL_TERMINAL_LINK {
								throughputDl = nc.ThroughputDl
							}
							if nc.ThroughputUl != DEFAULT_THROUGHPUT_UL_TERMINAL_LINK {
								throughputUl = nc.ThroughputUl
							}
							if nc.PacketLoss != DEFAULT_PACKET_LOSS_TERMINAL_LINK {
								packetLoss = nc.PacketLoss
							}
						}
						nl.NetChar = createNetChar(latency, latencyVariation, "", throughputDl, throughputUl, packetLoss)

						// Reset deprecated values to omit them
						nl.TerminalLinkLatency = 0
						nl.TerminalLinkLatencyVariation = 0
						nl.TerminalLinkPacketLoss = 0
						nl.TerminalLinkThroughput = 0

						// Physical Locations
						for iPl := range nl.PhysicalLocations {
							pl := &nl.PhysicalLocations[iPl]
							latency = pl.LinkLatency
							latencyVariation = pl.LinkLatencyVariation
							throughputDl = pl.LinkThroughput
							throughputUl = pl.LinkThroughput
							packetLoss = pl.LinkPacketLoss
							nc := pl.NetChar
							if nc != nil {
								//netchar got already created, if values are default, replace them with what was in the scenario, otherwise, leave them as is
								if nc.Latency != DEFAULT_LATENCY_LINK {
									latency = nc.Latency
								}
								if nc.LatencyVariation != DEFAULT_LATENCY_JITTER_LINK {
									latencyVariation = nc.LatencyVariation
								}
								if nc.ThroughputDl != DEFAULT_THROUGHPUT_DL_LINK {
									throughputDl = nc.ThroughputDl
								}
								if nc.ThroughputUl != DEFAULT_THROUGHPUT_UL_LINK {
									throughputUl = nc.ThroughputUl
								}
								if nc.PacketLoss != DEFAULT_PACKET_LOSS_LINK {
									packetLoss = nc.PacketLoss
								}
							}
							pl.NetChar = createNetChar(latency, latencyVariation, "", throughputDl, throughputUl, packetLoss)

							// Reset deprecated values to omit them
							pl.LinkLatency = 0
							pl.LinkLatencyVariation = 0
							pl.LinkPacketLoss = 0
							pl.LinkThroughput = 0

							for iProc := range pl.Processes {
								proc := &pl.Processes[iProc]
								latency = proc.AppLatency
								latencyVariation = proc.AppLatencyVariation
								throughputDl = proc.AppThroughput
								throughputUl = proc.AppThroughput
								packetLoss = proc.AppPacketLoss
								nc := proc.NetChar
								if nc != nil {
									//netchar got already created, if values are default, replace them with what was in the scenario, otherwise, leave them as is
									if nc.Latency != DEFAULT_LATENCY_APP {
										latency = nc.Latency
									}
									if nc.LatencyVariation != DEFAULT_LATENCY_JITTER_APP {
										latencyVariation = nc.LatencyVariation
									}
									if nc.ThroughputDl != DEFAULT_THROUGHPUT_DL_APP {
										throughputDl = nc.ThroughputDl
									}
									if nc.ThroughputUl != DEFAULT_THROUGHPUT_UL_APP {
										throughputUl = nc.ThroughputUl
									}
									if nc.PacketLoss != DEFAULT_PACKET_LOSS_APP {
										packetLoss = nc.PacketLoss
									}
								}
								proc.NetChar = createNetChar(latency, latencyVariation, "", throughputDl, throughputUl, packetLoss)

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
