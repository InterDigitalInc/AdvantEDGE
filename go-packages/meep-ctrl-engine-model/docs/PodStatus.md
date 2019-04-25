# PodStatus

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Pod name | [optional] [default to null]
**Namespace** | **string** | Pod namespace | [optional] [default to null]
**MeepApp** | **string** | Pod process name | [optional] [default to null]
**MeepOrigin** | **string** | Pod origin(core, scenario) | [optional] [default to null]
**MeepScenario** | **string** | Pod scenario name | [optional] [default to null]
**Phase** | **string** | Pod phase | [optional] [default to null]
**PodInitialized** | **string** | Pod initialized (true/false) | [optional] [default to null]
**PodReady** | **string** | Pod ready (true/false) | [optional] [default to null]
**PodScheduled** | **string** | Pod scheduled (true/false) | [optional] [default to null]
**PodUnschedulable** | **string** | Pod unschedulable (true/false) | [optional] [default to null]
**PodConditionError** | **string** | Pod error message | [optional] [default to null]
**ContainerStatusesMsg** | **string** | Failed container error message | [optional] [default to null]
**NbOkContainers** | **string** | Number of containers that are up | [optional] [default to null]
**NbTotalContainers** | **string** | Number of total containers in the pod | [optional] [default to null]
**NbPodRestart** | **string** | Number of container failures leading to pod restarts | [optional] [default to null]
**LogicalState** | **string** | State that is mapping the kubernetes api state | [optional] [default to null]
**StartTime** | **string** | Pod creation time | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


