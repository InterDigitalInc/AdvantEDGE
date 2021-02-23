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

	// Int32 Range
	err = validateInt32Range(0, 1, 10)
	if err == nil {
		t.Fatalf("Int32 should be out of range")
	}
	err = validateInt32Range(100, 1, 10)
	if err == nil {
		t.Fatalf("Int32 should be out of range")
	}
	err = validateInt32Range(1, 1, 10)
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = validateInt32Range(10, 1, 10)
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = validateInt32Range(5, 1, 10)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Float32 Range
	err = validateFloat32Range(0.1, 1.0, 10.0)
	if err == nil {
		t.Fatalf("Int32 should be out of range")
	}
	err = validateFloat32Range(100.0, 1.0, 10.0)
	if err == nil {
		t.Fatalf("Int32 should be out of range")
	}
	err = validateFloat32Range(1.0, 1.0, 10.0)
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = validateFloat32Range(10.0, 1.0, 10.0)
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = validateFloat32Range(5.0, 1.0, 10.0)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Float64 Range
	err = validateFloat64Range(0.1, 1.0, 10.0)
	if err == nil {
		t.Fatalf("Int32 should be out of range")
	}
	err = validateFloat64Range(100.0, 1.0, 10.0)
	if err == nil {
		t.Fatalf("Int32 should be out of range")
	}
	err = validateFloat64Range(1.0, 1.0, 10.0)
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = validateFloat64Range(10.0, 1.0, 10.0)
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = validateFloat64Range(5.0, 1.0, 10.0)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Enum
	enumVals := []string{"val1", "val2", "val3"}
	err = validateStringEnum("", enumVals)
	if err == nil {
		t.Fatalf("String should not be in enum")
	}
	err = validateStringEnum("dummy", enumVals)
	if err == nil {
		t.Fatalf("String should not be in enum")
	}
	err = validateStringEnum("val1", enumVals)
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = validateStringEnum("val2", enumVals)
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = validateStringEnum("val3", enumVals)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// MAC Address
	err = validateMacAddress("badmac")
	if err == nil {
		t.Fatalf("MAC Address should be invalid")
	}
	err = validateMacAddress("11:22:33")
	if err == nil {
		t.Fatalf("MAC Address should be invalid")
	}
	err = validateMacAddress("11:22:33:44:55:66")
	if err == nil {
		t.Fatalf("MAC Address should be invalid")
	}
	err = validateMacAddress("")
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = validateMacAddress("0123456789ab")
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = validateMacAddress("CDEF01234567")
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Wireless Type List
	err = validateWirelessTypeList("badtype")
	if err == nil {
		t.Fatalf("Wireless Type should be invalid")
	}
	err = validateWirelessTypeList("3g")
	if err == nil {
		t.Fatalf("Wireless Type should be invalid")
	}
	err = validateWirelessTypeList("4g,none")
	if err == nil {
		t.Fatalf("Wireless Type should be invalid")
	}
	err = validateWirelessTypeList("wifi")
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = validateWirelessTypeList("5g")
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = validateWirelessTypeList("4g")
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = validateWirelessTypeList("other")
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = validateWirelessTypeList("wifi,5g,4g,other")
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = validateWirelessTypeList("wifi, 4g,   5g")
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Path
	err = validatePath("", true)
	if err == nil {
		t.Fatalf("Path should be present")
	}
	err = validatePath("invalid path", true)
	if err == nil {
		t.Fatalf("Path should be invalid")
	}
	err = validatePath("", false)
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = validatePath("/valid/path", true)
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = validatePath("another/valid/path", true)
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = validatePath("registry:repo/name", true)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Environment Variables
	err = validateEnvVar("INVALID VAR=value")
	if err == nil {
		t.Fatalf("Env var format should be invalid")
	}
	err = validateEnvVar("VAR=value,INVALID_VAR")
	if err == nil {
		t.Fatalf("Env var format should be invalid")
	}
	err = validateEnvVar("")
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = validateEnvVar("VAR=!@#$%^&*(),VAR2='val with spaces,VAR_NO_VAL=")
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestValidatePhyLoc(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// pl := dataModel.PhysicalLocation{Id: "", Name: "my-pl"}
	// err := validatePhyLoc(&pl)
	// if err != nil {
	// 	t.Fatalf(err.Error())
	// }

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

	// pl := dataModel.Process{Id: "", Name: "my-proc"}
	// err := validateProc(&pl)
	// if err != nil {
	// 	t.Fatalf(err.Error())
	// }

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
