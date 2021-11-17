# AdvantEdgeMecApplicationSupportApi.TimingCapsNtpServers

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ntpServerAddrType** | **String** | Address type of NTP server | 
**ntpServerAddr** | **String** | NTP server address | 
**minPollingInterval** | **Number** | Minimum poll interval for NTP messages, in seconds as a power of two. Range 3...17 | 
**maxPollingInterval** | **Number** | Maximum poll interval for NTP messages, in seconds as a power of two. Range 3...17 | 
**localPriority** | **Number** | NTP server local priority | 
**authenticationOption** | **String** | NTP authentication option | 
**authenticationKeyNum** | **Number** | Authentication key number | 


<a name="NtpServerAddrTypeEnum"></a>
## Enum: NtpServerAddrTypeEnum


* `IP_ADDRESS` (value: `"IP_ADDRESS"`)

* `DNS_NAME` (value: `"DNS_NAME"`)




<a name="AuthenticationOptionEnum"></a>
## Enum: AuthenticationOptionEnum


* `NONE` (value: `"NONE"`)

* `SYMMETRIC_KEY` (value: `"SYMMETRIC_KEY"`)

* `AUTO_KEY` (value: `"AUTO_KEY"`)




