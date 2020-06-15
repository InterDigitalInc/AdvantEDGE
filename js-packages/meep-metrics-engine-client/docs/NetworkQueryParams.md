# AdvantEdgeMetricsServiceRestApi.NetworkQueryParams

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**tags** | [**[Tag]**](Tag.md) | Tag names to match in query. Supported values:<br> <li>src: Source network element name <li>dest: Destination network element name | [optional] 
**fields** | **[String]** | Field names to return in query response. Supported values:<br> <li>lat: Round-trip latency (ms)<br> <li>ul: Uplink throughput from src to dest (Mbps) <li>dl: Downlink throughput from dest to src (Mbps) <li>ulos: Uplink packet loss from src to dest (%) <li>dlos: Downlink packet loss from dest to src (%) | [optional] 
**scope** | [**Scope**](Scope.md) |  | [optional] 


<a name="[FieldsEnum]"></a>
## Enum: [FieldsEnum]


* `lat` (value: `"lat"`)

* `ul` (value: `"ul"`)

* `dl` (value: `"dl"`)

* `ulos` (value: `"ulos"`)

* `dlos` (value: `"dlos"`)




