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
	"fmt"
	"regexp"
	"strconv"
	"strings"

	dataModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-model"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"

	"github.com/blang/semver"
	"github.com/google/uuid"
)

// Validator status types
const (
	ValidatorStatusValid   = "SCENARIO-VALID"
	ValidatorStatusUpdated = "SCENARIO-UPDATED"
	ValidatorStatusError   = "SCENARIO-ERROR"
)

const (
	REGEX_NAME               = `^(([a-z0-9][-a-z0-9.]*)?[a-z0-9])+$`
	REGEX_VARIABLE_NAME      = `^(([_a-z0-9A-Z][_-a-z0-9A-Z.]*)?[_a-z0-9A-Z])+$`
	REGEX_MAC_ADDRESS        = `^(([_a-f0-9A-F][_-a-f0-9A-Z]*)?[_a-f0-9A-F])+$`
	REGEX_WIRELESS_TYPE_LIST = `^((,\s*)?(wifi|5g|4g|other))+$`
	REGEX_PATH               = `[\^#%&$\*<>\?\{\|\} ]+`
	REGEX_DNN                = `^(([a-z0-9A-Z][-a-z0-9A-Z.]*)?[a-z0-9A-Z])+$`
	REGEX_ECSP               = `^(([a-z0-9A-Z][ a-z0-9A-Z]*)?[a-z0-9A-Z])+$`
)

const (
	LATENCY_MIN                  = 0
	LATENCY_MAX                  = 250000
	LATENCY_DISTRIBUTION_DEFAULT = "Normal"
	JITTER_MIN                   = 0
	JITTER_MAX                   = 250000
	PACKET_LOSS_MIN              = float64(0.0)
	PACKET_LOSS_MAX              = float64(100.0)
	THROUGHPUT_MIN               = 1
	THROUGHPUT_MAX               = 1000000
	THROUGHPUT_DEFAULT           = 1000
	VELOCITY_MIN                 = 0
	VELOCITY_MAX                 = 1000000
	RADIUS_MIN                   = 1
	RADIUS_MAX                   = 1000000
	SERVICE_PORT_MIN             = 1
	SERVICE_PORT_MAX             = 65535
	SERVICE_NODE_PORT_MIN        = 30000
	SERVICE_NODE_PORT_MAX        = 32767
	GPU_COUNT_MIN                = 1
	GPU_COUNT_MAX                = 4
	MIN_CPU_COUNT_MIN            = 0.1
	MIN_CPU_COUNT_MAX            = 100.0
	MAX_CPU_COUNT_MIN            = 0.1
	MAX_CPU_COUNT_MAX            = 100.0
	MIN_MEMORY_MIN               = 1
	MIN_MEMORY_MAX               = 1000000
	MAX_MEMORY_MIN               = 1
	MAX_MEMORY_MAX               = 1000000
)

// Enums
var LATENCY_DIST_ENUM = []string{"Normal", "Pareto", "Paretonormal", "Uniform"}
var EOP_MODE_ENUM = []string{"LOOP", "REVERSE"}
var GPU_TYPE_ENUM = []string{"NVIDIA"}
var PROTOCOL_ENUM = []string{"UDP", "TCP"}
var CONNECTIVITY_MODEL_ENUM = []string{"OPEN", "PDU"}

// Current validator version
var ValidatorVersion = semver.Version{Major: 1, Minor: 8, Patch: 0}

// Versions requiring scenario update
var Version130 = semver.Version{Major: 1, Minor: 3, Patch: 0}
var Version140 = semver.Version{Major: 1, Minor: 4, Patch: 0}
var Version150 = semver.Version{Major: 1, Minor: 5, Patch: 0}
var Version151 = semver.Version{Major: 1, Minor: 5, Patch: 1}
var Version153 = semver.Version{Major: 1, Minor: 5, Patch: 3}
var Version168 = semver.Version{Major: 1, Minor: 6, Patch: 8}

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
func ValidateScenario(jsonScenario []byte, name string) (validJsonScenario []byte, status string, err error) {
	var scenarioVersion semver.Version
	var scenarioUpdated = false

	// Unmarshal scenario
	scenario := new(dataModel.Scenario)
	err = json.Unmarshal(jsonScenario, scenario)
	if err != nil {
		log.Error(err.Error())
		return nil, ValidatorStatusError, err
	}

	if name != "" {
		if scenario.Name != name {
			err = errors.New("Scenario creation name " + name + " incompatible with scenario body content name " + scenario.Name + ". They must be the same.")
			return nil, ValidatorStatusError, err
		}
	}

	// Retrieve scenario version
	// If no version found, assume & set current validator version
	if scenario.Version == "" {
		scenarioVersion = ValidatorVersion
		scenario.Version = ValidatorVersion.String()
		scenarioUpdated = true
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
	}

	// Run upgrade functions starting from oldest applicable patch to newest

	// UPGRADE TO 1.4.0
	if scenarioVersion.LT(Version140) {
		upgradeScenarioTo140(scenario)
		scenarioVersion = Version140
		scenarioUpdated = true
	}
	// UPGRADE TO 1.5.0
	if scenarioVersion.LT(Version150) {
		upgradeScenarioTo150(scenario)
		scenarioVersion = Version150
		scenarioUpdated = true
	}
	// UPGRADE TO 1.5.1
	if scenarioVersion.LT(Version151) {
		upgradeScenarioTo151(scenario)
		scenarioVersion = Version151
		scenarioUpdated = true
	}
	// UPGRADE TO 1.5.3
	if scenarioVersion.LT(Version153) {
		upgradeScenarioTo153(scenario)
		scenarioVersion = Version153
		scenarioUpdated = true
	}
	// UPGRADE TO 1.6.8
	if scenarioVersion.LT(Version168) {
		upgradeScenarioTo168(scenario)
		scenarioVersion = Version168
		scenarioUpdated = true
	}

	// Set current scenario version
	if scenarioVersion.LT(ValidatorVersion) {
		scenario.Version = ValidatorVersion.String()
		scenarioUpdated = true
	}

	// Validate scenario format & content
	err = validateScenario(scenario)
	if err != nil {
		log.Error("Scenario validation failed for: " + scenario.Name)
		log.Error(err.Error())
		return nil, ValidatorStatusError, err
	}

	// Marshal updated scenario
	if scenarioUpdated {
		validJsonScenario, err = json.Marshal(scenario)
		if err != nil {
			return nil, ValidatorStatusError, err
		}
		return validJsonScenario, ValidatorStatusUpdated, err
	} else {
		return jsonScenario, ValidatorStatusValid, nil
	}
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
				LATENCY_DISTRIBUTION_DEFAULT,
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

func upgradeScenarioTo151(scenario *dataModel.Scenario) {
	//changes in 160 (151 for now) vs 150
	//rename POA-CELLULAR to POA-4G

	// Set updated version
	scenario.Version = Version151.String()

	// Migrate netchar information
	if scenario.Deployment != nil {
		for iDomain := range scenario.Deployment.Domains {
			domain := &scenario.Deployment.Domains[iDomain]
			for iZone := range domain.Zones {
				zone := &domain.Zones[iZone]
				for iNl := range zone.NetworkLocations {
					nl := &zone.NetworkLocations[iNl]
					if nl.Type_ == "POA-CELLULAR" {
						nl.Type_ = "POA-4G"
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

func upgradeScenarioTo153(scenario *dataModel.Scenario) {
	// Set updated version
	scenario.Version = Version153.String()

	// Set Physical location connection parameters
	if scenario.Deployment != nil {
		for iDomain := range scenario.Deployment.Domains {
			domain := &scenario.Deployment.Domains[iDomain]
			for iZone := range domain.Zones {
				zone := &domain.Zones[iZone]
				for iNl := range zone.NetworkLocations {
					nl := &zone.NetworkLocations[iNl]
					for iPl := range nl.PhysicalLocations {
						pl := &nl.PhysicalLocations[iPl]
						pl.Connected = true
						pl.WirelessType = ""
						if pl.Type_ == "UE" {
							pl.Wireless = true
						} else {
							pl.Wireless = false
						}
					}
				}
			}
		}
	}
}

func upgradeScenarioTo168(scenario *dataModel.Scenario) {
	// Set updated version
	scenario.Version = Version168.String()

	// Set default Connectivity Model
	if scenario.Deployment != nil {
		if scenario.Deployment.Connectivity == nil {
			scenario.Deployment.Connectivity = new(dataModel.ConnectivityConfig)
		}
		if scenario.Deployment.Connectivity.Model == "" {
			scenario.Deployment.Connectivity.Model = "OPEN"
		}
	}
}

// Validate scenario
func validateScenario(scenario *dataModel.Scenario) error {
	idMap := make(map[string]bool)
	nameMap := make(map[string]bool)

	// Validate scenario
	if scenario == nil {
		return errors.New("scenario == nil")
	}
	if err := validateUniqueId(scenario.Id, idMap); err != nil {
		return err
	}
	if err := validateUniqueName(scenario.Name, nameMap); err != nil {
		return err
	}

	// Validate deployment
	deployment := scenario.Deployment
	if err := validateDeployment(deployment); err != nil {
		return err
	}

	// Validate domains
	for domainIndex := range deployment.Domains {
		domain := &deployment.Domains[domainIndex]
		// TODO -- Add Domain validation
		if err := validateUniqueId(domain.Id, idMap); err != nil {
			return err
		}
		if err := validateUniqueName(domain.Name, nameMap); err != nil {
			return err
		}

		// Validate zones
		for zoneIndex := range domain.Zones {
			zone := &domain.Zones[zoneIndex]
			// TODO -- Add Zone validation
			if err := validateUniqueId(zone.Id, idMap); err != nil {
				return err
			}
			if err := validateUniqueName(zone.Name, nameMap); err != nil {
				return err
			}

			// Validate Network Locations
			for nlIndex := range zone.NetworkLocations {
				nl := &zone.NetworkLocations[nlIndex]
				// TODO -- Add NetworkLocation validation
				if err := validateUniqueId(nl.Id, idMap); err != nil {
					return err
				}
				if err := validateUniqueName(nl.Name, nameMap); err != nil {
					return err
				}

				// Validate Physical Locations
				for plIndex := range nl.PhysicalLocations {
					pl := &nl.PhysicalLocations[plIndex]
					if err := validatePhyLoc(pl); err != nil {
						return err
					}
					if err := validateUniqueId(pl.Id, idMap); err != nil {
						return err
					}
					if err := validateUniqueName(pl.Name, nameMap); err != nil {
						return err
					}

					// Validate Processes
					for procIndex := range pl.Processes {
						proc := &pl.Processes[procIndex]
						if err := validateProc(proc); err != nil {
							return err
						}
						if err := validateUniqueId(proc.Id, idMap); err != nil {
							return err
						}
						if err := validateUniqueName(proc.Name, nameMap); err != nil {
							return err
						}
					}
				}
			}
		}
	}
	return nil
}

func validateUniqueId(id string, idMap map[string]bool) error {
	if _, found := idMap[id]; found {
		return errors.New("Id not unique: " + id)
	}
	idMap[id] = true
	return nil
}

func validateUniqueName(name string, nameMap map[string]bool) error {
	if _, found := nameMap[name]; found {
		return errors.New("Name not unique: " + name)
	}
	nameMap[name] = true
	return nil
}

// Validate the Deployment
func validateDeployment(deployment *dataModel.Deployment) (err error) {
	// Deployment
	if deployment == nil {
		return errors.New("deployment == nil")
	}
	// Connectivity
	if deployment.Connectivity == nil {
		return errors.New("connectivity == nil")
	}
	// Connectivity Model
	err = validateStringEnum(deployment.Connectivity.Model, CONNECTIVITY_MODEL_ENUM)
	if err != nil {
		return errors.New("Invalid connectivity model: " + err.Error())
	}
	return nil
}

// Validate the provided Physical Location
func validatePhyLoc(pl *dataModel.PhysicalLocation) (err error) {
	// ID: Create new UUID if none provided
	if pl.Id == "" {
		pl.Id = uuid.New().String()
	}
	// Name
	err = validateName(pl.Name)
	if err != nil {
		return err
	}
	// Type
	if !IsPhyLoc(pl.Type_) {
		return errors.New("Unsupported PhysicalLocation Type: " + pl.Type_)
	}
	// MAC Address
	err = validateMacAddress(pl.MacId)
	if err != nil {
		return err
	}
	// Wireless Types
	err = validateWirelessTypeList(pl.WirelessType)
	if err != nil {
		return err
	}
	// Network Characteristics
	err = validateNetChar(pl.NetChar)
	if err != nil {
		return err
	}
	// DataNetwork
	err = validateDataNetwork(pl.DataNetwork, pl.Type_)
	if err != nil {
		return err
	}
	// GeoData
	err = validateGeoData(pl.GeoData, pl.Type_)
	if err != nil {
		return err
	}
	return nil
}

// Validate the provided Process
func validateProc(proc *dataModel.Process) (err error) {
	// ID: Create new UUID if none provided
	if proc.Id == "" {
		proc.Id = uuid.New().String()
	}
	// Name
	err = validateName(proc.Name)
	if err != nil {
		return err
	}
	// Type
	if !IsProc(proc.Type_) {
		return errors.New("Unsupported Process Type: " + proc.Type_)
	}
	// Network Characteristics
	err = validateNetChar(proc.NetChar)
	if err != nil {
		return err
	}

	// App deployment type-specific validation
	if proc.IsExternal {
		// EXTERNAL APP

		// Service Config
		err = validateExternalConfig(proc.ExternalConfig)
		if err != nil {
			return err
		}

		// TODO - Validate placement identifier

	} else if proc.UserChartLocation != "" {
		// USER-DEFINED CHART APP

		// User Chart Location
		err = validatePath(proc.UserChartLocation, true)
		if err != nil {
			return errors.New("Invalid user chart location: " + err.Error())
		}
		// User Chart Group
		err = validateChartGroup(proc.UserChartGroup)
		if err != nil {
			return err
		}
		// User Chart Alternate values
		err = validatePath(proc.UserChartAlternateValues, false)
		if err != nil {
			return errors.New("Invalid user chart alternate values: " + err.Error())
		}

	} else {
		// INTERNAL APP

		// Container Image Name
		err = validatePath(proc.Image, true)
		if err != nil {
			return errors.New("Invalid container image name: " + err.Error())
		}
		// GPU Config
		err = validateGpuConfig(proc.GpuConfig)
		if err != nil {
			return err
		}
		// CPU Config
		err = validateCpuConfig(proc.CpuConfig)
		if err != nil {
			return err
		}
		// Memory Config
		err = validateMemoryConfig(proc.MemoryConfig)
		if err != nil {
			return err
		}
		// Environment variables
		err = validateEnvVar(proc.Environment)
		if err != nil {
			return errors.New("Invalid env var: " + err.Error())
		}
		// Command
		err = validatePath(proc.CommandExe, false)
		if err != nil {
			return errors.New("Invalid command: " + err.Error())
		}
		// TODO - Validate command arguments
		// TODO - Validate placement identifier

		if proc.Type_ == NodeTypeEdgeApp || proc.Type_ == NodeTypeCloudApp {
			// Service Config
			err = validateServiceConfig(proc.ServiceConfig)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func validateNetChar(nc *dataModel.NetworkCharacteristics) (err error) {
	// Mandatory field
	if nc == nil {
		return errors.New("Missing network characteristics")
	}

	// Set default values
	if nc.ThroughputUl == 0 {
		nc.ThroughputUl = THROUGHPUT_DEFAULT
	}
	if nc.ThroughputDl == 0 {
		nc.ThroughputDl = THROUGHPUT_DEFAULT
	}
	if nc.LatencyDistribution == "" {
		nc.LatencyDistribution = LATENCY_DISTRIBUTION_DEFAULT
	}

	// Validate fields
	err = validateInt32Range(nc.Latency, LATENCY_MIN, LATENCY_MAX)
	if err != nil {
		return errors.New("Invalid latency: " + err.Error())
	}
	err = validateInt32Range(nc.LatencyVariation, JITTER_MIN, JITTER_MAX)
	if err != nil {
		return errors.New("Invalid jitter: " + err.Error())
	}
	err = validateStringEnum(nc.LatencyDistribution, LATENCY_DIST_ENUM)
	if err != nil {
		return errors.New("Invalid latency distribution: " + err.Error())
	}
	err = validateFloat64Range(nc.PacketLoss, PACKET_LOSS_MIN, PACKET_LOSS_MAX)
	if err != nil {
		return errors.New("Invalid packet loss: " + err.Error())
	}
	err = validateInt32Range(nc.ThroughputUl, THROUGHPUT_MIN, THROUGHPUT_MAX)
	if err != nil {
		return errors.New("Invalid UL throughput: " + err.Error())
	}
	err = validateInt32Range(nc.ThroughputDl, THROUGHPUT_MIN, THROUGHPUT_MAX)
	if err != nil {
		return errors.New("Invalid DL throughput: " + err.Error())
	}
	return nil
}

func validateDataNetwork(dn *dataModel.DnConfig, typ string) (err error) {
	// Optional field
	if dn == nil {
		return nil
	}
	// DNN
	err = validateDnn(dn.Dnn)
	if err != nil {
		return err
	}
	if typ == NodeTypeUE && dn.Dnn != "" {
		return errors.New("UE must not have a configured DNN")
	}
	// ECSP
	err = validateEcsp(dn.Ecsp)
	if err != nil {
		return err
	}
	if typ == NodeTypeUE && dn.Ecsp != "" {
		return errors.New("UE must not have a configured ECSP")
	}
	return nil
}

func validateGeoData(gd *dataModel.GeoData, typ string) (err error) {
	// Optional field
	if gd == nil {
		return nil
	}
	// Radius
	if IsNetLoc(typ) {
		err = validateFloat32Range(gd.Radius, RADIUS_MIN, RADIUS_MAX)
		if err != nil {
			return errors.New("Invalid radius: " + err.Error())
		}
	}
	// Path (optional)
	if gd.Path != nil {
		// End-of-path mode
		err = validateStringEnum(gd.EopMode, EOP_MODE_ENUM)
		if err != nil {
			return errors.New("Invalid EOP mode: " + err.Error())
		}
		// Velocity
		err = validateFloat32Range(gd.Velocity, VELOCITY_MIN, VELOCITY_MAX)
		if err != nil {
			return errors.New("Invalid velocity: " + err.Error())
		}
		// TODO - Validate Path
	}

	// TODO - Validate Location

	return nil
}

func validateGpuConfig(cfg *dataModel.GpuConfig) (err error) {
	// Optional field
	if cfg == nil {
		return nil
	}
	// GPU Count
	err = validateInt32Range(cfg.Count, GPU_COUNT_MIN, GPU_COUNT_MAX)
	if err != nil {
		return errors.New("Invalid GPU count: " + err.Error())
	}
	// GPU Type
	err = validateStringEnum(cfg.Type_, GPU_TYPE_ENUM)
	if err != nil {
		return errors.New("Invalid GPU type: " + err.Error())
	}
	return nil
}

func validateCpuConfig(cfg *dataModel.CpuConfig) (err error) {
	// Optional field
	if cfg == nil {
		return nil
	}
	// Min CPU Count (optional)
	if cfg.Min != 0 {
		err = validateFloat32Range(cfg.Min, MIN_CPU_COUNT_MIN, MIN_CPU_COUNT_MAX)
		if err != nil {
			return errors.New("Invalid min CPU count: " + err.Error())
		}
	}
	// Max CPU Count (optional)
	if cfg.Max != 0 {
		err = validateFloat32Range(cfg.Max, MAX_CPU_COUNT_MIN, MAX_CPU_COUNT_MAX)
		if err != nil {
			return errors.New("Invalid max CPU count: " + err.Error())
		}
	}
	// Max CPU must be greater than Min CPU
	if cfg.Max != 0 && cfg.Max < cfg.Min {
		return fmt.Errorf("Max CPU [%f] less than Min CPU [%f]", cfg.Max, cfg.Min)
	}
	return nil
}

func validateMemoryConfig(cfg *dataModel.MemoryConfig) (err error) {
	// Optional field
	if cfg == nil {
		return nil
	}
	// Min Memory (optional)
	if cfg.Min != 0 {
		err = validateInt32Range(cfg.Min, MIN_MEMORY_MIN, MIN_MEMORY_MAX)
		if err != nil {
			return errors.New("Invalid min memory: " + err.Error())
		}
	}
	// Max Memory (optional)
	if cfg.Max != 0 {
		err = validateInt32Range(cfg.Max, MAX_MEMORY_MIN, MAX_MEMORY_MAX)
		if err != nil {
			return errors.New("Invalid max memory: " + err.Error())
		}
	}
	// Max Memory must be greater than Min Memory
	if cfg.Max != 0 && cfg.Max < cfg.Min {
		return fmt.Errorf("Max Memory [%d] less than Min Memory [%d]", cfg.Max, cfg.Min)
	}
	return nil
}

func validateExternalConfig(cfg *dataModel.ExternalConfig) (err error) {
	// Ingress Service Mapping
	for _, svc := range cfg.IngressServiceMap {
		err = validateIngressSvc(&svc)
		if err != nil {
			return err
		}
	}
	// EgressServiceMapping
	for _, svc := range cfg.EgressServiceMap {
		err = validateEgressSvc(&svc)
		if err != nil {
			return err
		}
	}
	return nil
}

func validateServiceConfig(cfg *dataModel.ServiceConfig) (err error) {
	// Optional field
	if cfg == nil {
		return nil
	}
	// Service Name -- TBD if it needs to be same as unique element name
	err = validateName(cfg.Name)
	if err != nil {
		return errors.New("Invalid service name: " + err.Error())
	}
	// Multi-Edge (Group) Service Name (optional)
	if cfg.MeSvcName != "" {
		err = validateName(cfg.MeSvcName)
		if err != nil {
			return errors.New("Invalid group svc name: " + err.Error())
		}
	}
	// Service Ports
	for _, port := range cfg.Ports {
		// Protocol
		err = validateStringEnum(port.Protocol, PROTOCOL_ENUM)
		if err != nil {
			return errors.New("Invalid service protocol: " + err.Error())
		}
		// Port
		err = validateInt32Range(port.Port, SERVICE_PORT_MIN, SERVICE_PORT_MAX)
		if err != nil {
			return errors.New("Invalid service port: " + err.Error())
		}
		// External Port
		if port.ExternalPort != 0 {
			err = validateInt32Range(port.ExternalPort, SERVICE_NODE_PORT_MIN, SERVICE_NODE_PORT_MAX)
			if err != nil {
				return errors.New("Invalid service ext port: " + err.Error())
			}
		}
	}
	return nil
}

func validateChartGroup(group string) (err error) {
	if group != "" {
		fields := strings.Split(group, ":")
		if len(fields) != 4 {
			return errors.New("Group format must be 'svc instance:svc group name:port:protocol'")
		}
		// Svc name
		err = validateFullName(fields[0])
		if err != nil {
			return errors.New("Invalid chart group svc name: " + err.Error())
		}
		// Svc group name (optional)
		if fields[1] != "" {
			err = validateName(fields[1])
			if err != nil {
				return errors.New("Invalid chart group name: " + err.Error())
			}
		}
		// Port
		port, err := strconv.Atoi(fields[2])
		if err != nil {
			return errors.New("Invalid chart group port: " + err.Error())
		}
		err = validateInt32Range(int32(port), SERVICE_PORT_MIN, SERVICE_PORT_MAX)
		if err != nil {
			return errors.New("Invalid chart group port: " + err.Error())
		}
		// Protocol
		err = validateStringEnum(fields[3], PROTOCOL_ENUM)
		if err != nil {
			return errors.New("Invalid chart group protocol: " + err.Error())
		}
	}
	return nil
}

func validateIngressSvc(svc *dataModel.IngressService) (err error) {
	// External Port
	err = validateInt32Range(svc.ExternalPort, SERVICE_NODE_PORT_MIN, SERVICE_NODE_PORT_MAX)
	if err != nil {
		return errors.New("Invalid Ingress ext port: " + err.Error())
	}
	// Svc name
	err = validateFullName(svc.Name)
	if err != nil {
		return errors.New("Invalid Ingress svc name: " + err.Error())
	}
	// Port
	err = validateInt32Range(svc.Port, SERVICE_PORT_MIN, SERVICE_PORT_MAX)
	if err != nil {
		return errors.New("Invalid Ingress port: " + err.Error())
	}
	// Protocol
	err = validateStringEnum(svc.Protocol, PROTOCOL_ENUM)
	if err != nil {
		return errors.New("Invalid Ingress protocol: " + err.Error())
	}
	return nil
}

func validateEgressSvc(svc *dataModel.EgressService) (err error) {
	// Svc name
	err = validateFullName(svc.Name)
	if err != nil {
		return errors.New("Invalid Egress svc name: " + err.Error())
	}
	// Group Svc name (optional)
	if svc.MeSvcName != "" {
		err = validateFullName(svc.MeSvcName)
		if err != nil {
			return errors.New("Invalid Egress group svc name: " + err.Error())
		}
	}
	// Port
	err = validateInt32Range(svc.Port, SERVICE_PORT_MIN, SERVICE_PORT_MAX)
	if err != nil {
		return errors.New("Invalid Egress port: " + err.Error())
	}
	// Protocol
	err = validateStringEnum(svc.Protocol, PROTOCOL_ENUM)
	if err != nil {
		return errors.New("Invalid Egress protocol: " + err.Error())
	}

	// TODO -- Validate IP

	return nil
}

func validateName(name string) (err error) {
	if name == "" {
		return errors.New("Name not provided")
	}
	if len(name) > 30 {
		return errors.New("Name length exceeds maximum of 30 characters")
	}
	matched, err := regexp.MatchString(REGEX_NAME, name)
	if err != nil || !matched {
		return errors.New("Name must be lowercase alphanumeric or '-' or '.'")
	}
	return nil
}

func validateFullName(name string) (err error) {
	if name == "" {
		return errors.New("Full name not provided")
	}
	if len(name) > 60 {
		return errors.New("Full name length exceeds maximum of 60 characters")
	}
	matched, err := regexp.MatchString(REGEX_NAME, name)
	if err != nil || !matched {
		return errors.New("Full name must be lowercase alphanumeric or '-' or '.'")
	}
	return nil
}

func validateVariableName(name string) (err error) {
	if name == "" {
		return errors.New("Variable name not provided")
	}
	if len(name) > 30 {
		return errors.New("Variable name length exceeds maximum of 30 characters")
	}
	matched, err := regexp.MatchString(REGEX_VARIABLE_NAME, name)
	if err != nil || !matched {
		return errors.New("Variable name must be alphanumeric or '-' or '.'")
	}
	return nil
}

func validateInt32Range(val int32, min int32, max int32) (err error) {
	if val < min || val > max {
		return fmt.Errorf("Int32 val: %d not in valid range: [%d - %d]", val, min, max)
	}
	return nil
}

func validateFloat32Range(val float32, min float32, max float32) (err error) {
	if val < min || val > max {
		return fmt.Errorf("Float32 val: %f not in valid range: [%f - %f]", val, min, max)
	}
	return nil
}

func validateFloat64Range(val float64, min float64, max float64) (err error) {
	if val < min || val > max {
		return fmt.Errorf("Float64 val: %f not in valid range: [%f - %f]", val, min, max)
	}
	return nil
}

func validateStringEnum(val string, enum []string) (err error) {
	for _, str := range enum {
		if val == str {
			return nil
		}
	}
	return fmt.Errorf("String val: %s not in enum: %v", val, enum)
}

func validateMacAddress(mac string) (err error) {
	if mac != "" {
		if len(mac) > 12 {
			return errors.New("MAC address > 12 chars")
		}
		matched, err := regexp.MatchString(REGEX_MAC_ADDRESS, mac)
		if err != nil || !matched {
			return errors.New("MAC address must be alphanumeric hex")
		}
	}
	return nil
}

func validateWirelessTypeList(list string) (err error) {
	if list != "" {
		matched, err := regexp.MatchString(REGEX_WIRELESS_TYPE_LIST, list)
		if err != nil || !matched {
			return errors.New("Wireless type list must be comma-separated list: wifi|5g|4g|other")
		}
	}
	return nil
}

func validateDnn(dnn string) (err error) {
	if dnn != "" {
		matched, err := regexp.MatchString(REGEX_DNN, dnn)
		if err != nil || !matched {
			return errors.New("DNN must be alphanumeric or '-' or '.'")
		}
	}
	return nil
}

func validateEcsp(ecsp string) (err error) {
	if ecsp != "" {
		matched, err := regexp.MatchString(REGEX_ECSP, ecsp)
		if err != nil || !matched {
			return errors.New("ECSP must be alphanumeric or ' '")
		}
	}
	return nil
}

func validatePath(path string, isRequired bool) (err error) {
	if path != "" {
		matched, err := regexp.MatchString(REGEX_PATH, path)
		if err != nil || matched {
			return fmt.Errorf("Invalid path format: %s", path)
		}
	} else if isRequired {
		return errors.New("Missing path")
	}
	return nil
}

func validateEnvVar(env string) (err error) {
	if env != "" {
		variables := strings.Split(env, ",")
		for _, variable := range variables {
			fields := strings.Split(variable, "=")
			if len(fields) != 2 {
				return fmt.Errorf("Invalid env var format: %s", variable)
			}
			err = validateVariableName(strings.TrimSpace(fields[0]))
			if err != nil {
				return err
			}
		}
	}
	return nil
}
