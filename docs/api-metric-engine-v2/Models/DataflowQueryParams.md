# DataflowQueryParams
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**tags** | [**List**](Tag.md) | Tag names to match in query. Supported values:&lt;br&gt;  | [optional] [default to null]
**scope** | [**Scope**](Scope.md) |  | [optional] [default to null]
**fields** | [**List**](string.md) | Requested information. Supported values:&lt;br&gt;&lt;li&gt;mermaid: Mermaid format&lt;br&gt; | [optional] [default to null]
**responseType** | [**String**](string.md) | Queried response Type. Supported Values:&lt;br&gt; NOTE1: only one of listonly or responly may be included  NOTE2: if listonly or responly are not included, the response contains both the list and string  &lt;li&gt;listonly: Include only a list of dataflow metrics in response&lt;br&gt; &lt;li&gt;stronly: Include only a concatenated string of dataflow metrics in response&lt;br&gt;  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

