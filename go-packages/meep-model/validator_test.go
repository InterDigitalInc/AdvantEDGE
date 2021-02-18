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
	"fmt"
	"testing"

	dataModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-model"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

func TestValidateComponents(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Name
	err := validateName("")
	if err == nil {
		t.Fatalf("Name with empty string should fail")
	}
	err = validateName("0123456789012345678901234567890")
	if err == nil {
		t.Fatalf("Name len > 30 should fail")
	}
	err = validateName("InvalidName")
	if err == nil {
		t.Fatalf("Name with invalid chars should fail")
	}
	err = validateName("invalid name")
	if err == nil {
		t.Fatalf("Name with invalid chars should fail")
	}
	err = validateName("my-valid.name123")
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = validateName("012345678901234567890123456789")
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Full Name
	err = validateFullName("")
	if err == nil {
		t.Fatalf("Full name with empty string should fail")
	}
	err = validateFullName("0123456789012345678901234567890123456789012345678901234567890")
	if err == nil {
		t.Fatalf("Full name len > 60 should fail")
	}
	err = validateFullName("InvalidName")
	if err == nil {
		t.Fatalf("Full name with invalid chars should fail")
	}
	err = validateFullName("invalid name")
	if err == nil {
		t.Fatalf("Full name with invalid chars should fail")
	}
	err = validateFullName("my-valid.name123")
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = validateFullName("012345678901234567890123456789012345678901234567890123456789")
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Variable Name
	err = validateVariableName("")
	if err == nil {
		t.Fatalf("Variable name with empty string should fail")
	}
	err = validateVariableName("0123456789012345678901234567890")
	if err == nil {
		t.Fatalf("Variable name len > 30 should fail")
	}
	err = validateVariableName("InvalidName*")
	if err == nil {
		t.Fatalf("Variable name with invalid chars should fail")
	}
	err = validateVariableName("invalid name")
	if err == nil {
		t.Fatalf("Variable name with invalid chars should fail")
	}
	err = validateVariableName("My-Valid.Name123")
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = validateVariableName("012345678901234567890123456789")
	if err != nil {
		t.Fatalf(err.Error())
	}

}

func TestValidatePhyLoc(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	pl := dataModel.PhysicalLocation{Id: "", Name: "my-pl"}
	err := validatePhyLoc(&pl)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Type_ string `json:"type,omitempty"`
	// // true: Physical location is external to MEEP false: Physical location is internal to MEEP
	// IsExternal              bool     `json:"isExternal,omitempty"`
	// GeoData                 *GeoData `json:"geoData,omitempty"`
	// NetworkLocationsInRange []string `json:"networkLocationsInRange,omitempty"`
	// // true: Physical location has network connectivity false: Physical location has no network connectivity
	// Connected bool `json:"connected,omitempty"`
	// // true: Physical location uses a wireless connection false: Physical location uses a wired connection
	// Wireless bool `json:"wireless,omitempty"`
	// // Prioritized, comma-separated list of supported wireless connection types. Default priority if not specififed is 'wifi,5g,4g,other'. Wireless connection types: - 4g - 5g - wifi - other
	// WirelessType string `json:"wirelessType,omitempty"`
	// // Key/Value Pair Map (string, string)
	// Meta map[string]string `json:"meta,omitempty"`
	// // Key/Value Pair Map (string, string)
	// UserMeta  map[string]string       `json:"userMeta,omitempty"`
	// Processes []Process               `json:"processes,omitempty"`
	// NetChar   *NetworkCharacteristics `json:"netChar,omitempty"`
	// // **DEPRECATED** As of release 1.5.0, replaced by netChar latency
	// LinkLatency int32 `json:"linkLatency,omitempty"`
	// // **DEPRECATED** As of release 1.5.0, replaced by netChar latencyVariation
	// LinkLatencyVariation int32 `json:"linkLatencyVariation,omitempty"`
	// // **DEPRECATED** As of release 1.5.0, replaced by netChar throughputUl and throughputDl
	// LinkThroughput int32 `json:"linkThroughput,omitempty"`
	// // **DEPRECATED** As of release 1.5.0, replaced by netChar packetLoss
	// LinkPacketLoss float64 `json:"linkPacketLoss,omitempty"`
	// // Physical location MAC Address
	// MacId string `json:"macId,omitempty"`

}

func TestValidateProc(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	pl := dataModel.Process{Id: "", Name: "my-proc"}
	err := validateProc(&pl)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// // true: process is external to MEEP false: process is internal to MEEP
	// IsExternal bool `json:"isExternal,omitempty"`
	// // Docker image to deploy inside MEEP
	// Image string `json:"image,omitempty"`
	// // Environment variables using the format NAME=\"value\",NAME=\"value\",NAME=\"value\"
	// Environment string `json:"environment,omitempty"`
	// // Arguments to command executable
	// CommandArguments string `json:"commandArguments,omitempty"`
	// // Executable to invoke at container start up
	// CommandExe     string          `json:"commandExe,omitempty"`
	// ServiceConfig  *ServiceConfig  `json:"serviceConfig,omitempty"`
	// GpuConfig      *GpuConfig      `json:"gpuConfig,omitempty"`
	// MemoryConfig   *MemoryConfig   `json:"memoryConfig,omitempty"`
	// CpuConfig      *CpuConfig      `json:"cpuConfig,omitempty"`
	// ExternalConfig *ExternalConfig `json:"externalConfig,omitempty"`
	// // Process status
	// Status string `json:"status,omitempty"`
	// // Chart location for the deployment of the chart provided by the user
	// UserChartLocation string `json:"userChartLocation,omitempty"`
	// // Chart values.yaml file location for the deployment of the chart provided by the user
	// UserChartAlternateValues string `json:"userChartAlternateValues,omitempty"`
	// // Chart supplemental information related to the group (service)
	// UserChartGroup string `json:"userChartGroup,omitempty"`
	// // Key/Value Pair Map (string, string)
	// Meta map[string]string `json:"meta,omitempty"`
	// // Key/Value Pair Map (string, string)
	// UserMeta map[string]string       `json:"userMeta,omitempty"`
	// NetChar  *NetworkCharacteristics `json:"netChar,omitempty"`
	// // Identifier used for process placement in AdvantEDGE cluster
	// PlacementId string `json:"placementId,omitempty"`
}
