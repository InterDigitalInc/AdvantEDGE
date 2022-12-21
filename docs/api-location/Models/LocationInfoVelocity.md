# LocationInfoVelocity
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**bearing** | [**Integer**](integer.md) | Bearing, expressed in the range 0° to 360°, as defined in ETSI TS 123 032 [14]. | [default to null]
**horizontalSpeed** | [**Integer**](integer.md) | Horizontal speed, expressed in km/h and defined in ETSI TS 123 032 [14]. | [default to null]
**uncertainty** | [**Integer**](integer.md) | Horizontal uncertainty, as defined in ETSI TS 123 032 [14]. Present only if \&quot;velocityType\&quot; equals 3 or 4 | [optional] [default to null]
**velocityType** | [**Integer**](integer.md) | Velocity information, as detailed in ETSI TS 123 032 [14], associated with the reported location coordinate: &lt;p&gt;1 &#x3D; HORIZONTAL &lt;p&gt;2 &#x3D; HORIZONTAL_VERTICAL &lt;p&gt;3 &#x3D; HORIZONTAL_UNCERT &lt;p&gt;4 &#x3D; HORIZONTAL_VERTICAL_UNCERT | [default to null]
**verticalSpeed** | [**Integer**](integer.md) | Vertical speed, expressed in km/h and defined in ETSI TS 123 032 [14]. Present only if \&quot;velocityType\&quot; equals 2 or 4 | [optional] [default to null]
**verticalUncertainty** | [**Integer**](integer.md) | Vertical uncertainty, as defined in ETSI TS 123 032 [14]. Present only if \&quot;velocityType\&quot; equals 4 | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

