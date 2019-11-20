# AdvantEdgePlatformControllerRestApi.PodStatus

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**name** | **String** | Pod name | [optional] 
**namespace** | **String** | Pod namespace | [optional] 
**meepApp** | **String** | Pod process name | [optional] 
**meepOrigin** | **String** | Pod origin(core, scenario) | [optional] 
**meepScenario** | **String** | Pod scenario name | [optional] 
**phase** | **String** | Pod phase | [optional] 
**podInitialized** | **String** | Pod initialized (true/false) | [optional] 
**podReady** | **String** | Pod ready (true/false) | [optional] 
**podScheduled** | **String** | Pod scheduled (true/false) | [optional] 
**podUnschedulable** | **String** | Pod unschedulable (true/false) | [optional] 
**podConditionError** | **String** | Pod error message | [optional] 
**containerStatusesMsg** | **String** | Failed container error message | [optional] 
**nbOkContainers** | **String** | Number of containers that are up | [optional] 
**nbTotalContainers** | **String** | Number of total containers in the pod | [optional] 
**nbPodRestart** | **String** | Number of container failures leading to pod restarts | [optional] 
**logicalState** | **String** | State that is mapping the kubernetes api state | [optional] 
**startTime** | **String** | Pod creation time | [optional] 


