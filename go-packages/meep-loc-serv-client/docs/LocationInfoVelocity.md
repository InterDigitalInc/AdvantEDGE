# LocationInfoVelocity

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Bearing** | **int32** | Bearing, expressed in the range 0° to 360°, as defined in [14]. | [default to null]
**HorizontalSpeed** | **int32** | Horizontal speed, expressed in km/h and defined in [14]. | [default to null]
**Uncertainty** | **int32** | Horizontal uncertainty, as defined in [14]. Present only if \&quot;velocityType\&quot; equals 3 or 4 | [optional] [default to null]
**VelocityType** | **int32** | Velocity information, as detailed in [14], associated with the reported location coordinate: &lt;p&gt;1 &#x3D; HORIZONTAL &lt;p&gt;2 &#x3D; HORIZONTAL_VERTICAL &lt;p&gt;3 &#x3D; HORIZONTAL_UNCERT &lt;p&gt;4 &#x3D; HORIZONTAL_VERTICAL_UNCERT | [default to null]
**VerticalSpeed** | **int32** | Vertical speed, expressed in km/h and defined in [14]. Present only if \&quot;velocityType\&quot; equals 2 or 4 | [optional] [default to null]
**VerticalUncertainty** | **int32** | Vertical uncertainty, as defined in [14]. Present only if \&quot;velocityType\&quot; equals 4 | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


