# InlineSubscription
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**\_links** | [**AssocStaSubscription__links**](AssocStaSubscription__links.md) |  | [optional] [default to null]
**apId** | [**ApIdentity**](ApIdentity.md) |  | [default to null]
**callbackReference** | [**URI**](URI.md) | URI selected by the service consumer to receive notifications on the subscribed WLAN Access Information Service. This shall be included both in the request and in response. | [default to null]
**expiryDeadline** | [**TimeStamp**](TimeStamp.md) |  | [optional] [default to null]
**subscriptionType** | [**String**](string.md) | Shall be set to \&quot;StaDataRateSubscription\&quot;. | [default to null]
**staId** | [**List**](StaIdentity.md) | Identifier(s) to uniquely specify the target client station(s) for the subscription | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

