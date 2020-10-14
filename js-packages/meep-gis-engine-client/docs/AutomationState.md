# AdvantEdgeGisEngineRestApi.AutomationState

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**type** | **String** | Automation type.<br> Automation loop evaluates enabled automation types once every second.<br> <p>Supported Types: <li>MOBILITY - Sends Mobility events to Sanbox Controller when UE changes POA. <li>MOVEMENT - Advances UEs along configured paths using previous position & velocity as inputs. <li>POAS-IN-RANGE - Sends POAS-IN-RANGE events to Sanbox Controller when list of POAs in range changes. <li>NETWORK-CHARACTERISTICS-UPDATE - Sends network characteristics update events to Sanbox Controller when throughput values change. | [optional] 
**active** | **Boolean** | Automation feature state | [optional] 


<a name="TypeEnum"></a>
## Enum: TypeEnum


* `MOBILITY` (value: `"MOBILITY"`)

* `MOVEMENT` (value: `"MOVEMENT"`)

* `POAS_IN_RANGE` (value: `"POAS-IN-RANGE"`)

* `NETWORK_CHARACTERISTICS_UPDATE` (value: `"NETWORK-CHARACTERISTICS-UPDATE"`)




