# Documentation for AdvantEDGE Sandbox Controller REST API

<a name="documentation-for-api-endpoints"></a>
## Documentation for API Endpoints

All URIs are relative to *http://localhost/sandbox-ctrl/v1*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*ActiveScenarioApi* | [**activateScenario**](Apis/ActiveScenarioApi.md#activatescenario) | **POST** /active/{name} | Deploy a scenario
*ActiveScenarioApi* | [**getActiveNodeServiceMaps**](Apis/ActiveScenarioApi.md#getactivenodeservicemaps) | **GET** /active/serviceMaps | Get deployed scenario's port mapping
*ActiveScenarioApi* | [**getActiveScenario**](Apis/ActiveScenarioApi.md#getactivescenario) | **GET** /active | Get the deployed scenario
*ActiveScenarioApi* | [**terminateScenario**](Apis/ActiveScenarioApi.md#terminatescenario) | **DELETE** /active | Terminate the deployed scenario
*EventReplayApi* | [**createReplayFile**](Apis/EventReplayApi.md#createreplayfile) | **POST** /replay/{name} | Add a replay file
*EventReplayApi* | [**createReplayFileFromScenarioExec**](Apis/EventReplayApi.md#createreplayfilefromscenarioexec) | **POST** /replay/{name}/generate | Generate a replay file from Active Scenario events
*EventReplayApi* | [**deleteReplayFile**](Apis/EventReplayApi.md#deletereplayfile) | **DELETE** /replay/{name} | Delete a replay file
*EventReplayApi* | [**deleteReplayFileList**](Apis/EventReplayApi.md#deletereplayfilelist) | **DELETE** /replay | Delete all replay files
*EventReplayApi* | [**getReplayFile**](Apis/EventReplayApi.md#getreplayfile) | **GET** /replay/{name} | Get a specific replay file
*EventReplayApi* | [**getReplayFileList**](Apis/EventReplayApi.md#getreplayfilelist) | **GET** /replay | Get all replay file names
*EventReplayApi* | [**getReplayStatus**](Apis/EventReplayApi.md#getreplaystatus) | **GET** /replaystatus | Get status of replay manager
*EventReplayApi* | [**loopReplay**](Apis/EventReplayApi.md#loopreplay) | **POST** /replay/{name}/loop | Loop-Execute a replay file present in the platform store
*EventReplayApi* | [**playReplayFile**](Apis/EventReplayApi.md#playreplayfile) | **POST** /replay/{name}/play | Execute a replay file present in the platform store
*EventReplayApi* | [**stopReplayFile**](Apis/EventReplayApi.md#stopreplayfile) | **POST** /replay/{name}/stop | Stop execution of a replay file
*EventsApi* | [**sendEvent**](Apis/EventsApi.md#sendevent) | **POST** /events/{type} | Send events to the deployed scenario


<a name="documentation-for-models"></a>
## Documentation for Models

 - [ActivationInfo](./Models/ActivationInfo.md)
 - [CellularDomainConfig](./Models/CellularDomainConfig.md)
 - [CellularPoaConfig](./Models/CellularPoaConfig.md)
 - [CpuConfig](./Models/CpuConfig.md)
 - [Deployment](./Models/Deployment.md)
 - [Domain](./Models/Domain.md)
 - [EgressService](./Models/EgressService.md)
 - [Event](./Models/Event.md)
 - [EventMobility](./Models/EventMobility.md)
 - [EventNetworkCharacteristicsUpdate](./Models/EventNetworkCharacteristicsUpdate.md)
 - [EventPoasInRange](./Models/EventPoasInRange.md)
 - [EventScenarioUpdate](./Models/EventScenarioUpdate.md)
 - [ExternalConfig](./Models/ExternalConfig.md)
 - [GeoData](./Models/GeoData.md)
 - [GpuConfig](./Models/GpuConfig.md)
 - [IngressService](./Models/IngressService.md)
 - [LineString](./Models/LineString.md)
 - [MemoryConfig](./Models/MemoryConfig.md)
 - [NetworkCharacteristics](./Models/NetworkCharacteristics.md)
 - [NetworkLocation](./Models/NetworkLocation.md)
 - [NodeDataUnion](./Models/NodeDataUnion.md)
 - [NodeServiceMaps](./Models/NodeServiceMaps.md)
 - [PhysicalLocation](./Models/PhysicalLocation.md)
 - [Poa4GConfig](./Models/Poa4GConfig.md)
 - [Poa5GConfig](./Models/Poa5GConfig.md)
 - [PoaWifiConfig](./Models/PoaWifiConfig.md)
 - [Point](./Models/Point.md)
 - [Process](./Models/Process.md)
 - [Replay](./Models/Replay.md)
 - [ReplayEvent](./Models/ReplayEvent.md)
 - [ReplayFileList](./Models/ReplayFileList.md)
 - [ReplayInfo](./Models/ReplayInfo.md)
 - [ReplayStatus](./Models/ReplayStatus.md)
 - [Scenario](./Models/Scenario.md)
 - [ScenarioConfig](./Models/ScenarioConfig.md)
 - [ScenarioNode](./Models/ScenarioNode.md)
 - [ServiceConfig](./Models/ServiceConfig.md)
 - [ServicePort](./Models/ServicePort.md)
 - [Zone](./Models/Zone.md)


<a name="documentation-for-authorization"></a>
## Documentation for Authorization

All endpoints do not require authorization.
