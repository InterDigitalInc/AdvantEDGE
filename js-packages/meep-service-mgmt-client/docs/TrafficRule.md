# AdvantEdgeMecApplicationSupportApi.TrafficRule

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**trafficRuleId** | **String** | Identify the traffic rule. | 
**filterType** | **String** | Definition of filter per FLOW or PACKET. If flow the filter match UE-&gt;EPC packet and the reverse packet is handled in the same context | 
**priority** | **Number** | Priority of this traffic rule. If traffic rule conflicts, the one with higher priority take precedence | 
**trafficFilter** | [**[TrafficFilter]**](TrafficFilter.md) |  | 
**action** | **String** | The action of the MEC host data plane when a packet matches the trafficFilter | 
**dstInterface** | [**DestinationInterface**](DestinationInterface.md) |  | [optional] 
**state** | **String** | Contains the traffic rule state. This attribute may be updated using HTTP PUT   method | 


<a name="FilterTypeEnum"></a>
## Enum: FilterTypeEnum


* `FLOW` (value: `"FLOW"`)

* `PACKET` (value: `"PACKET"`)




<a name="ActionEnum"></a>
## Enum: ActionEnum


* `DROP` (value: `"DROP"`)

* `FORWARD_DECAPSULATED` (value: `"FORWARD_DECAPSULATED"`)

* `FORWARD_ENCAPSULATED` (value: `"FORWARD_ENCAPSULATED"`)

* `PASSTHROUGH` (value: `"PASSTHROUGH"`)

* `DUPLICATE_DECAPSULATED` (value: `"DUPLICATE_DECAPSULATED"`)

* `DUPLICATE_ENCAPSULATED` (value: `"DUPLICATE_ENCAPSULATED"`)




<a name="StateEnum"></a>
## Enum: StateEnum


* `ACTIVE` (value: `"ACTIVE"`)

* `INACTIVE` (value: `"INACTIVE"`)




