/*
 * Copyright (c) 2022  The AdvantEDGE Authors
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

func TestValidateName(t *testing.T) {
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
}

func TestValidateFullName(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Full Name
	err := validateFullName("")
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
}

func TestValidateVariableName(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Variable Name
	err := validateVariableName("")
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

func TestValidateInt32Range(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Int32 Range
	err := validateInt32Range(0, 1, 10)
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
}

func TestValidateFloat32Range(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Float32 Range
	err := validateFloat32Range(0.1, 1.0, 10.0)
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
}

func TestValidateFloat64Range(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Float64 Range
	err := validateFloat64Range(0.1, 1.0, 10.0)
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
}

func TestValidateStringEnum(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Enum
	enumVals := []string{"val1", "val2", "val3"}
	err := validateStringEnum("", enumVals)
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
}

func TestValidateMacAddress(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// MAC Address
	err := validateMacAddress("badmac")
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
}

func TestValidateWirelessTypeList(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Wireless Type List
	err := validateWirelessTypeList("badtype")
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
	err = validateWirelessTypeList("d2d")
	if err != nil {
		t.Fatalf(err.Error())
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
	err = validateWirelessTypeList("d2d,wifi,5g,4g,other")
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = validateWirelessTypeList("d2d,   wifi, 4g,   5g")
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestValidatePath(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Path
	err := validatePath("", true)
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
}

func TestValidateEnvVar(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Environment Variables
	err := validateEnvVar("INVALID VAR=value")
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

func TestValidateIngressSvc(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Ingress Service
	svc := dataModel.IngressService{Name: "", Port: 1, ExternalPort: 30000, Protocol: "TCP"}
	err := validateIngressSvc(&svc)
	if err == nil {
		t.Fatalf("Ingress should be invalid")
	}
	svc = dataModel.IngressService{Name: "svc", Port: 0, ExternalPort: 30000, Protocol: "UDP"}
	err = validateIngressSvc(&svc)
	if err == nil {
		t.Fatalf("Ingress should be invalid")
	}
	svc = dataModel.IngressService{Name: "svc", Port: 32767, ExternalPort: 1000, Protocol: "TCP"}
	err = validateIngressSvc(&svc)
	if err == nil {
		t.Fatalf("Ingress should be invalid")
	}
	svc = dataModel.IngressService{Name: "svc", Port: 1234, ExternalPort: 30000, Protocol: "invalid"}
	err = validateIngressSvc(&svc)
	if err == nil {
		t.Fatalf("Ingress should be invalid")
	}
	svc = dataModel.IngressService{Name: "svc", Port: 1000, ExternalPort: 31000, Protocol: "TCP"}
	err = validateIngressSvc(&svc)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestValidateEgressSvc(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Egress Service
	svc := dataModel.EgressService{Name: "", MeSvcName: "me-svc", Port: 1, Protocol: "TCP"}
	err := validateEgressSvc(&svc)
	if err == nil {
		t.Fatalf("Egress should be invalid")
	}
	svc = dataModel.EgressService{Name: "svc", MeSvcName: "me-svc", Port: 0, Protocol: "TCP"}
	err = validateEgressSvc(&svc)
	if err == nil {
		t.Fatalf("Egress should be invalid")
	}
	svc = dataModel.EgressService{Name: "svc", MeSvcName: "me-svc", Port: 1234, Protocol: "invalid"}
	err = validateEgressSvc(&svc)
	if err == nil {
		t.Fatalf("Egress should be invalid")
	}
	svc = dataModel.EgressService{Name: "svc", MeSvcName: "", Port: 32767, Protocol: "UDP"}
	err = validateEgressSvc(&svc)
	if err != nil {
		t.Fatalf(err.Error())
	}
	svc = dataModel.EgressService{Name: "svc", MeSvcName: "me-svc", Port: 1000, Protocol: "TCP"}
	err = validateEgressSvc(&svc)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestValidateChartGroup(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Chart Group
	err := validateChartGroup("invalid field count:2nd arg")
	if err == nil {
		t.Fatalf("Chart group should be invalid")
	}
	err = validateChartGroup(":me-svc:1000:TCP")
	if err == nil {
		t.Fatalf("Chart group should be invalid")
	}
	err = validateChartGroup("svc:me-svc::TCP")
	if err == nil {
		t.Fatalf("Chart group should be invalid")
	}
	err = validateChartGroup("svc:me-svc:1000:")
	if err == nil {
		t.Fatalf("Chart group should be invalid")
	}
	err = validateChartGroup("")
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = validateChartGroup("svc::1000:TCP")
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = validateChartGroup("svc:me-svc:1000:UDP")
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestValidateServiceConfig(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Service Config
	// port1 := ServicePort
	ports := []dataModel.ServicePort{}
	cfg := dataModel.ServiceConfig{Name: "", MeSvcName: "me-svc", Ports: ports}
	err := validateServiceConfig(&cfg)
	if err == nil {
		t.Fatalf("Service Config should be invalid")
	}
	ports = []dataModel.ServicePort{dataModel.ServicePort{Protocol: "invalid", Port: 0, ExternalPort: 0}}
	cfg = dataModel.ServiceConfig{Name: "svc", MeSvcName: "me-svc", Ports: ports}
	err = validateServiceConfig(&cfg)
	if err == nil {
		t.Fatalf("Service Config should be invalid")
	}
	ports = []dataModel.ServicePort{dataModel.ServicePort{Protocol: "UDP", Port: 1111, ExternalPort: 1000}}
	cfg = dataModel.ServiceConfig{Name: "svc", MeSvcName: "me-svc", Ports: ports}
	err = validateServiceConfig(&cfg)
	if err == nil {
		t.Fatalf("Service Config should be invalid")
	}
	ports = []dataModel.ServicePort{
		dataModel.ServicePort{Protocol: "UDP", Port: 1111, ExternalPort: 0},
		dataModel.ServicePort{Protocol: "TCP", Port: 0, ExternalPort: 0},
	}
	cfg = dataModel.ServiceConfig{Name: "svc", MeSvcName: "me-svc", Ports: ports}
	err = validateServiceConfig(&cfg)
	if err == nil {
		t.Fatalf("Service Config should be invalid")
	}
	err = validateServiceConfig(nil)
	if err != nil {
		t.Fatalf(err.Error())
	}
	ports = []dataModel.ServicePort{}
	cfg = dataModel.ServiceConfig{Name: "svc", MeSvcName: "", Ports: ports}
	err = validateServiceConfig(&cfg)
	if err != nil {
		t.Fatalf(err.Error())
	}
	ports = []dataModel.ServicePort{{Protocol: "UDP", Port: 1111, ExternalPort: 31000}}
	cfg = dataModel.ServiceConfig{Name: "svc", MeSvcName: "", Ports: ports}
	err = validateServiceConfig(&cfg)
	if err != nil {
		t.Fatalf(err.Error())
	}
	ports = []dataModel.ServicePort{
		dataModel.ServicePort{Protocol: "UDP", Port: 1111, ExternalPort: 0},
		dataModel.ServicePort{Protocol: "TCP", Port: 2222, ExternalPort: 31000},
	}
	cfg = dataModel.ServiceConfig{Name: "svc", MeSvcName: "", Ports: ports}
	err = validateServiceConfig(&cfg)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestValidateGpuConfig(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// GPU Config
	cfg := dataModel.GpuConfig{Count: 0, Type_: "NVIDIA"}
	err := validateGpuConfig(&cfg)
	if err == nil {
		t.Fatalf("GPU Config should be invalid")
	}
	cfg = dataModel.GpuConfig{Count: 5, Type_: "NVIDIA"}
	err = validateGpuConfig(&cfg)
	if err == nil {
		t.Fatalf("GPU Config should be invalid")
	}
	cfg = dataModel.GpuConfig{Count: 1, Type_: "invalid"}
	err = validateGpuConfig(&cfg)
	if err == nil {
		t.Fatalf("GPU Config should be invalid")
	}
	err = validateGpuConfig(nil)
	if err != nil {
		t.Fatalf(err.Error())
	}
	cfg = dataModel.GpuConfig{Count: 2, Type_: "NVIDIA"}
	err = validateGpuConfig(&cfg)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestValidateCpuConfig(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// GPU Config
	cfg := dataModel.CpuConfig{Min: 1001, Max: 10}
	err := validateCpuConfig(&cfg)
	if err == nil {
		t.Fatalf("CPU Config should be invalid")
	}
	cfg = dataModel.CpuConfig{Min: 0, Max: -1}
	err = validateCpuConfig(&cfg)
	if err == nil {
		t.Fatalf("CPU Config should be invalid")
	}
	cfg = dataModel.CpuConfig{Min: 2, Max: 1}
	err = validateCpuConfig(&cfg)
	if err == nil {
		t.Fatalf("CPU Config should be invalid")
	}
	err = validateCpuConfig(nil)
	if err != nil {
		t.Fatalf(err.Error())
	}
	cfg = dataModel.CpuConfig{Min: 0, Max: 100}
	err = validateCpuConfig(&cfg)
	if err != nil {
		t.Fatalf(err.Error())
	}
	cfg = dataModel.CpuConfig{Min: 0.1, Max: 0}
	err = validateCpuConfig(&cfg)
	if err != nil {
		t.Fatalf(err.Error())
	}
	cfg = dataModel.CpuConfig{Min: 1, Max: 2}
	err = validateCpuConfig(&cfg)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestValidateMemoryConfig(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// GPU Config
	cfg := dataModel.MemoryConfig{Min: 1000001, Max: 10}
	err := validateMemoryConfig(&cfg)
	if err == nil {
		t.Fatalf("Memory Config should be invalid")
	}
	cfg = dataModel.MemoryConfig{Min: 0, Max: -1}
	err = validateMemoryConfig(&cfg)
	if err == nil {
		t.Fatalf("Memory Config should be invalid")
	}
	cfg = dataModel.MemoryConfig{Min: 1000, Max: 500}
	err = validateMemoryConfig(&cfg)
	if err == nil {
		t.Fatalf("Memory Config should be invalid")
	}
	err = validateMemoryConfig(nil)
	if err != nil {
		t.Fatalf(err.Error())
	}
	cfg = dataModel.MemoryConfig{Min: 0, Max: 1000000}
	err = validateMemoryConfig(&cfg)
	if err != nil {
		t.Fatalf(err.Error())
	}
	cfg = dataModel.MemoryConfig{Min: 1, Max: 0}
	err = validateMemoryConfig(&cfg)
	if err != nil {
		t.Fatalf(err.Error())
	}
	cfg = dataModel.MemoryConfig{Min: 250, Max: 1000}
	err = validateMemoryConfig(&cfg)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestValidatePhyLoc(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// TODO
}

func TestValidateProc(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// TODO
}
