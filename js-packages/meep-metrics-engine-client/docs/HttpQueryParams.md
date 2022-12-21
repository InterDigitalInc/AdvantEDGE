# AdvantEdgeMetricsServiceRestApi.HttpQueryParams

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**tags** | [**[Tag]**](Tag.md) | Tag names to match in query. Supported values:<br> <li>logger_name: Logger instances that issued the http notification or processed the request <li>msg_type: Http metric type (request, response, notification) | [optional] 
**fields** | **[String]** | Field names to return in query response. Supported values:<br> <li>id: Http metrics identifier<br> <li>endpoint: Http metrics queried endpoint<br> <li>url: Http metrics queried endpoint with query parameters<br> <li>method: Http metrics method<br> <li>resp_code: Http metrics response status code<br> <li>resp_body: Http metrics response body<br> <li>body: Http metrics body<br> <li>proc_time: Request processing time in ms<br> <li>logger_name: Logger instances that issued the http notification or processed the request<br> <li>msg_type: Http metric type (request, response, notification)<br> <li>direction: DEPRECATED -- Http metric direction (RX, TX) | [optional] 
**scope** | [**Scope**](Scope.md) |  | [optional] 


<a name="[FieldsEnum]"></a>
## Enum: [FieldsEnum]


* `id` (value: `"id"`)

* `endpoint` (value: `"endpoint"`)

* `url` (value: `"url"`)

* `method` (value: `"method"`)

* `respCode` (value: `"resp_code"`)

* `respBody` (value: `"resp_body"`)

* `body` (value: `"body"`)

* `procTime` (value: `"proc_time"`)

* `loggerName` (value: `"logger_name"`)

* `msgType` (value: `"msg_type"`)

* `direction` (value: `"direction"`)




