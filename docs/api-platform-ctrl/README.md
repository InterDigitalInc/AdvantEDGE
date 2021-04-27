# Documentation for AdvantEDGE Platform Controller REST API

<a name="documentation-for-api-endpoints"></a>
## Documentation for API Endpoints

All URIs are relative to *http://localhost/platform-ctrl/v1*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*SandboxControlApi* | [**createSandbox**](Apis/SandboxControlApi.md#createsandbox) | **POST** /sandboxes | Create a new sandbox
*SandboxControlApi* | [**createSandboxWithName**](Apis/SandboxControlApi.md#createsandboxwithname) | **POST** /sandboxes/{name} | Create a new sandbox
*SandboxControlApi* | [**deleteSandbox**](Apis/SandboxControlApi.md#deletesandbox) | **DELETE** /sandboxes/{name} | Delete a specific sandbox
*SandboxControlApi* | [**deleteSandboxList**](Apis/SandboxControlApi.md#deletesandboxlist) | **DELETE** /sandboxes | Delete all active sandboxes
*SandboxControlApi* | [**getSandbox**](Apis/SandboxControlApi.md#getsandbox) | **GET** /sandboxes/{name} | Get a specific sandbox
*SandboxControlApi* | [**getSandboxList**](Apis/SandboxControlApi.md#getsandboxlist) | **GET** /sandboxes | Get all active sandboxes
*ScenarioConfigurationApi* | [**createScenario**](Apis/ScenarioConfigurationApi.md#createscenario) | **POST** /scenarios/{name} | Add a scenario
*ScenarioConfigurationApi* | [**deleteScenario**](Apis/ScenarioConfigurationApi.md#deletescenario) | **DELETE** /scenarios/{name} | Delete a scenario
*ScenarioConfigurationApi* | [**deleteScenarioList**](Apis/ScenarioConfigurationApi.md#deletescenariolist) | **DELETE** /scenarios | Delete all scenarios
*ScenarioConfigurationApi* | [**getScenario**](Apis/ScenarioConfigurationApi.md#getscenario) | **GET** /scenarios/{name} | Get a specific scenario
*ScenarioConfigurationApi* | [**getScenarioList**](Apis/ScenarioConfigurationApi.md#getscenariolist) | **GET** /scenarios | Get all scenarios
*ScenarioConfigurationApi* | [**setScenario**](Apis/ScenarioConfigurationApi.md#setscenario) | **PUT** /scenarios/{name} | Update a scenario


<a name="documentation-for-models"></a>
## Documentation for Models

 - [CellularDomainConfig](./Models/CellularDomainConfig.md)
 - [CellularPoaConfig](./Models/CellularPoaConfig.md)
 - [ConnectivityConfig](./Models/ConnectivityConfig.md)
 - [CpuConfig](./Models/CpuConfig.md)
 - [DNConfig](./Models/DNConfig.md)
 - [Deployment](./Models/Deployment.md)
 - [Domain](./Models/Domain.md)
 - [EgressService](./Models/EgressService.md)
 - [ExternalConfig](./Models/ExternalConfig.md)
 - [GeoData](./Models/GeoData.md)
 - [GpuConfig](./Models/GpuConfig.md)
 - [IngressService](./Models/IngressService.md)
 - [LineString](./Models/LineString.md)
 - [MemoryConfig](./Models/MemoryConfig.md)
 - [NetworkCharacteristics](./Models/NetworkCharacteristics.md)
 - [NetworkLocation](./Models/NetworkLocation.md)
 - [PhysicalLocation](./Models/PhysicalLocation.md)
 - [Poa4GConfig](./Models/Poa4GConfig.md)
 - [Poa5GConfig](./Models/Poa5GConfig.md)
 - [PoaWifiConfig](./Models/PoaWifiConfig.md)
 - [Point](./Models/Point.md)
 - [Process](./Models/Process.md)
 - [Sandbox](./Models/Sandbox.md)
 - [SandboxConfig](./Models/SandboxConfig.md)
 - [SandboxList](./Models/SandboxList.md)
 - [Scenario](./Models/Scenario.md)
 - [ScenarioConfig](./Models/ScenarioConfig.md)
 - [ScenarioList](./Models/ScenarioList.md)
 - [ServiceConfig](./Models/ServiceConfig.md)
 - [ServicePort](./Models/ServicePort.md)
 - [Zone](./Models/Zone.md)


<a name="documentation-for-authorization"></a>
## Documentation for Authorization

All endpoints do not require authorization.
