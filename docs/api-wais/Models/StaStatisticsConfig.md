# StaStatisticsConfig
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**groupIdentity** | [**Integer**](integer.md) | As per Table 9-92 of IEEE 802.11-2016 [8]. | [default to null]
**measurementCount** | [**Integer**](integer.md) | Valid if triggeredReport &#x3D; true. Specifies the number of MAC service data units or protocol data units to determine if the trigger conditions are met. | [optional] [default to null]
**triggerCondition** | [**STACounterTriggerCondition**](STACounterTriggerCondition.md) |  | [optional] [default to null]
**triggerTimeout** | [**Integer**](integer.md) | Valid if triggeredReport &#x3D; true. The Trigger Timeout field contains a value in units of 100 time-units of 1 024 µs during which a measuring STA does not generate further triggered STA Statistics Reports after a trigger condition has been met. | [optional] [default to null]
**triggeredReport** | [**Boolean**](boolean.md) | True &#x3D; triggered reporting, otherwise duration. | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

