# HttpQueryParams
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**tags** | [**List**](Tag.md) | Tag names to match in query. Supported values:&lt;br&gt; &lt;li&gt;logger_name: Logger instances that issued the http notification or processed the request &lt;li&gt;direction: Notification or Request type of http metric | [optional] [default to null]
**fields** | [**List**](string.md) | Field names to return in query response. Supported values:&lt;br&gt; &lt;li&gt;id: Http metrics identifier&lt;br&gt; &lt;li&gt;endpoint: Http metrics queried endpoint&lt;br&gt; &lt;li&gt;url: Http metrics queried endpoint with query parameters&lt;br&gt; &lt;li&gt;method: Http metrics method&lt;br&gt; &lt;li&gt;resp_code: Http metrics response status code&lt;br&gt; &lt;li&gt;resp_body: Http metrics response body&lt;br&gt; &lt;li&gt;body: Http metrics body&lt;br&gt; &lt;li&gt;proc_time: Request processing time in ms | [optional] [default to null]
**scope** | [**Scope**](Scope.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

