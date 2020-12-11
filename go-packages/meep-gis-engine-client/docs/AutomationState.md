# AutomationState

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Type_** | **string** | Automation type.&lt;br&gt; Automation loop evaluates enabled automation types once every second.&lt;br&gt; &lt;p&gt;Supported Types: &lt;li&gt;MOBILITY - Sends Mobility events to Sanbox Controller when UE changes POA. &lt;li&gt;MOVEMENT - Advances UEs along configured paths using previous position &amp; velocity as inputs. &lt;li&gt;POAS-IN-RANGE - Sends POAS-IN-RANGE events to Sanbox Controller when list of POAs in range changes. &lt;li&gt;NETWORK-CHARACTERISTICS-UPDATE - Sends network characteristics update events to Sanbox Controller when throughput values change. | [optional] [default to null]
**Active** | **bool** | Automation feature state | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


