# AdvantEdgeMetricsServiceRestApi.SeqQueryParams

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**tags** | [**[Tag]**](Tag.md) | Tag names to match in query. Supported values:<br>  | [optional] 
**fields** | **[String]** | Requested information. Supported values:<br> NOTE: only one of mermaid or sdorg must be included  <li>mermaid: Mermaid format<br> <li>sdorg: Sequencediagram.org format<br>  | [optional] 
**responseType** | **String** | Queried response Type. Supported Values:<br> NOTE1: only one of listonly or responly may be included  NOTE2: if listonly or responly are not included, the response contains both the list and string  <li>listonly: Include only a list of sequence metrics in response<br> <li>stronly: Include only a concatenated string of sequence metrics in response<br>  | [optional] 
**scope** | [**Scope**](Scope.md) |  | [optional] 


<a name="[FieldsEnum]"></a>
## Enum: [FieldsEnum]


* `mermaid` (value: `"mermaid"`)

* `sdorg` (value: `"sdorg"`)




<a name="ResponseTypeEnum"></a>
## Enum: ResponseTypeEnum


* `listonly` (value: `"listonly"`)

* `stronly` (value: `"stronly"`)




