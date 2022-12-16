# Documentation for AdvantEDGE MEC Application Support API

<a name="documentation-for-api-endpoints"></a>
## Documentation for API Endpoints

All URIs are relative to *https://localhost/sandboxname/mec_app_support/v1*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*MecAppSupportApi* | [**applicationsConfirmReadyPOST**](Apis/MecAppSupportApi.md#applicationsconfirmreadypost) | **POST** /applications/{appInstanceId}/confirm_ready | This method may be used by the MEC application instance to notify the MEC platform that it is up and running. 
*MecAppSupportApi* | [**applicationsConfirmTerminationPOST**](Apis/MecAppSupportApi.md#applicationsconfirmterminationpost) | **POST** /applications/{appInstanceId}/confirm_termination | This method is used to confirm the application level termination  of an application instance.
*MecAppSupportApi* | [**applicationsSubscriptionDELETE**](Apis/MecAppSupportApi.md#applicationssubscriptiondelete) | **DELETE** /applications/{appInstanceId}/subscriptions/{subscriptionId} | This method deletes a mecAppSuptApiSubscription. This method is typically used in \"Unsubscribing from service availability event notifications\" procedure.
*MecAppSupportApi* | [**applicationsSubscriptionGET**](Apis/MecAppSupportApi.md#applicationssubscriptionget) | **GET** /applications/{appInstanceId}/subscriptions/{subscriptionId} | The GET method requests information about a subscription for this requestor. Upon success, the response contains entity body with the subscription for the requestor.
*MecAppSupportApi* | [**applicationsSubscriptionsGET**](Apis/MecAppSupportApi.md#applicationssubscriptionsget) | **GET** /applications/{appInstanceId}/subscriptions | The GET method may be used to request information about all subscriptions for this requestor. Upon success, the response contains entity body with all the subscriptions for the requestor.
*MecAppSupportApi* | [**applicationsSubscriptionsPOST**](Apis/MecAppSupportApi.md#applicationssubscriptionspost) | **POST** /applications/{appInstanceId}/subscriptions | The POST method may be used to create a new subscription. One example use case is to create a new subscription to the MEC service availability notifications. Upon success, the response contains entity body describing the created subscription.
*MecAppSupportApi* | [**timingCapsGET**](Apis/MecAppSupportApi.md#timingcapsget) | **GET** /timing/timing_caps | This method retrieves the information of the platform's timing capabilities which corresponds to the timing capabilities query
*MecAppSupportApi* | [**timingCurrentTimeGET**](Apis/MecAppSupportApi.md#timingcurrenttimeget) | **GET** /timing/current_time | This method retrieves the information of the platform's current time which corresponds to the get platform time procedure
*UnsupportedApi* | [**applicationsDnsRuleGET**](Apis/UnsupportedApi.md#applicationsdnsruleget) | **GET** /applications/{appInstanceId}/dns_rules/{dnsRuleId} | This method retrieves information about a DNS rule associated with a MEC application instance.
*UnsupportedApi* | [**applicationsDnsRulePUT**](Apis/UnsupportedApi.md#applicationsdnsruleput) | **PUT** /applications/{appInstanceId}/dns_rules/{dnsRuleId} | This method activates, de-activates or updates a traffic rule.
*UnsupportedApi* | [**applicationsDnsRulesGET**](Apis/UnsupportedApi.md#applicationsdnsrulesget) | **GET** /applications/{appInstanceId}/dns_rules | This method retrieves information about all the DNS rules associated with a MEC application instance.
*UnsupportedApi* | [**applicationsTrafficRuleGET**](Apis/UnsupportedApi.md#applicationstrafficruleget) | **GET** /applications/{appInstanceId}/traffic_rules/{trafficRuleId} | This method retrieves information about all the traffic rules associated with a MEC application instance.
*UnsupportedApi* | [**applicationsTrafficRulePUT**](Apis/UnsupportedApi.md#applicationstrafficruleput) | **PUT** /applications/{appInstanceId}/traffic_rules/{trafficRuleId} | This method retrieves information about all the traffic rules associated with a MEC application instance.
*UnsupportedApi* | [**applicationsTrafficRulesGET**](Apis/UnsupportedApi.md#applicationstrafficrulesget) | **GET** /applications/{appInstanceId}/traffic_rules | This method retrieves information about all the traffic rules associated with a MEC application instance.


<a name="documentation-for-models"></a>
## Documentation for Models

 - [AppReadyConfirmation](./Models/AppReadyConfirmation.md)
 - [AppTerminationConfirmation](./Models/AppTerminationConfirmation.md)
 - [AppTerminationNotification](./Models/AppTerminationNotification.md)
 - [AppTerminationNotificationLinks](./Models/AppTerminationNotificationLinks.md)
 - [AppTerminationNotificationSubscription](./Models/AppTerminationNotificationSubscription.md)
 - [CurrentTime](./Models/CurrentTime.md)
 - [DestinationInterface](./Models/DestinationInterface.md)
 - [DestinationInterfaceInterfaceType](./Models/DestinationInterfaceInterfaceType.md)
 - [DnsRule](./Models/DnsRule.md)
 - [DnsRuleIpAddressType](./Models/DnsRuleIpAddressType.md)
 - [DnsRuleState](./Models/DnsRuleState.md)
 - [LinkType](./Models/LinkType.md)
 - [LinkTypeConfirmTermination](./Models/LinkTypeConfirmTermination.md)
 - [MecAppSuptApiSubscriptionLinkList](./Models/MecAppSuptApiSubscriptionLinkList.md)
 - [MecAppSuptApiSubscriptionLinkListLinks](./Models/MecAppSuptApiSubscriptionLinkListLinks.md)
 - [MecAppSuptApiSubscriptionLinkListSubscription](./Models/MecAppSuptApiSubscriptionLinkListSubscription.md)
 - [OperationActionType](./Models/OperationActionType.md)
 - [ProblemDetails](./Models/ProblemDetails.md)
 - [Self](./Models/Self.md)
 - [TimeSourceStatus](./Models/TimeSourceStatus.md)
 - [TimingCaps](./Models/TimingCaps.md)
 - [TimingCapsNtpServers](./Models/TimingCapsNtpServers.md)
 - [TimingCapsNtpServersAuthenticationOption](./Models/TimingCapsNtpServersAuthenticationOption.md)
 - [TimingCapsNtpServersNtpServerAddrType](./Models/TimingCapsNtpServersNtpServerAddrType.md)
 - [TimingCapsPtpMasters](./Models/TimingCapsPtpMasters.md)
 - [TimingCapsTimeStamp](./Models/TimingCapsTimeStamp.md)
 - [TrafficFilter](./Models/TrafficFilter.md)
 - [TrafficRule](./Models/TrafficRule.md)
 - [TrafficRuleAction](./Models/TrafficRuleAction.md)
 - [TrafficRuleFilterType](./Models/TrafficRuleFilterType.md)
 - [TrafficRuleState](./Models/TrafficRuleState.md)
 - [TunnelInfo](./Models/TunnelInfo.md)
 - [TunnelInfoTunnelType](./Models/TunnelInfoTunnelType.md)


<a name="documentation-for-authorization"></a>
## Documentation for Authorization

All endpoints do not require authorization.
