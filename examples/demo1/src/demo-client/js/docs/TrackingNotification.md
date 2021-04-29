# MeepDemoAppApi.TrackingNotification

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**callbackData** | **String** | CallBackData if passed by the application during the associated Subscription (Zone or User Tracking) operation | 
**zoneId** | **String** | Unique Identifier of a Location Zone | [optional] 
**address** | **String** | Address of the user or device based on the connected access point - address = acr:<UE IP address> | [optional] 
**interestRealm** | **String** | Details about the access point, geographical position, industry, etc. | [optional] 
**userEventType** | [**UserEventType**](UserEventType.md) |  | [optional] 
**currentAccessPointId** | **String** | Unique identifier of a point of access | [optional] 
**previousAccessPointId** | **String** | Unique identifier of a point of access | [optional] 
**timestamp** | [**TimeStamp**](TimeStamp.md) |  | [optional] 


