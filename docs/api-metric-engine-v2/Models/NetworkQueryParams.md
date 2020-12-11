# NetworkQueryParams
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**tags** | [**List**](Tag.md) | Tag names to match in query. Supported values:&lt;br&gt; &lt;li&gt;src: Source network element name &lt;li&gt;dest: Destination network element name | [optional] [default to null]
**fields** | [**List**](string.md) | Field names to return in query response. Supported values:&lt;br&gt; &lt;li&gt;lat: Round-trip latency (ms)&lt;br&gt; &lt;li&gt;ul: Uplink throughput from src to dest (Mbps) &lt;li&gt;dl: Downlink throughput from dest to src (Mbps) &lt;li&gt;ulos: Uplink packet loss from src to dest (%) &lt;li&gt;dlos: Downlink packet loss from dest to src (%) | [optional] [default to null]
**scope** | [**Scope**](Scope.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

