# AdvantEdgeMecApplicationSupportApi.DnsRule

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**dnsRuleId** | **String** | Identifies the DNS Rule | 
**domainName** | **String** | FQDN resolved by the DNS rule | 
**ipAddressType** | **String** | IP address type | 
**ipAddress** | **String** | IP address associated with the FQDN resolved by the DNS rule | 
**ttl** | **Number** | Time to live value | [optional] 
**state** | **String** | DNS rule state. This attribute may be updated using HTTP PUT method | 


<a name="IpAddressTypeEnum"></a>
## Enum: IpAddressTypeEnum


* `V6` (value: `"IP_V6"`)

* `V4` (value: `"IP_V4"`)




<a name="StateEnum"></a>
## Enum: StateEnum


* `ACTIVE` (value: `"ACTIVE"`)

* `INACTIVE` (value: `"INACTIVE"`)




