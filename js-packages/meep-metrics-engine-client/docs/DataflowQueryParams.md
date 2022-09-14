# AdvantEdgeMetricsServiceRestApi.DataflowQueryParams

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**tags** | [**[Tag]**](Tag.md) | Tag names to match in query. Supported values:<br>  | [optional] 
**scope** | [**Scope**](Scope.md) |  | [optional] 
**fields** | **[String]** | Requested information. Supported values:<br><li>mermaid: Mermaid format<br> | [optional] 
**responseType** | **String** | Queried response Type. Supported Values:<br> NOTE1: only one of listonly or responly may be included  NOTE2: if listonly or responly are not included, the response contains both the list and string  <li>listonly: Include only a list of dataflow metrics in response<br> <li>stronly: Include only a concatenated string of dataflow metrics in response<br>  | [optional] 


<a name="[FieldsEnum]"></a>
## Enum: [FieldsEnum]


* `mermaid` (value: `"mermaid"`)




<a name="ResponseTypeEnum"></a>
## Enum: ResponseTypeEnum


* `listonly` (value: `"listonly"`)

* `stronly` (value: `"stronly"`)




