# LocationInfo

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Accuracy** | **int32** | Horizontal accuracy / (semi-major) uncertainty of location provided in meters, as defined in [14]. Present only if \&quot;shape\&quot; equals 4, 5 or 6 | [optional] [default to null]
**AccuracyAltitude** | **int32** | Altitude accuracy / uncertainty of location provided in meters, as defined in [14]. Present only if \&quot;shape\&quot; equals 3 or 4 | [optional] [default to null]
**AccuracySemiMinor** | **int32** | Horizontal accuracy / (semi-major) uncertainty of location provided in meters, as defined in [14]. Present only if \&quot;shape\&quot; equals 4, 5 or 6 | [optional] [default to null]
**Altitude** | **float32** | Location altitude relative to the WGS84 ellipsoid surface. | [optional] [default to null]
**Confidence** | **int32** | Confidence by which the position of a target entity is known to be within the shape description, expressed as a percentage and defined in [14]. Present only if \&quot;shape\&quot; equals 1, 4 or 6 | [optional] [default to null]
**IncludedAngle** | **int32** | Present only if \&quot;shape\&quot; equals 6 | [optional] [default to null]
**InnerRadius** | **int32** | Present only if \&quot;shape\&quot; equals 6 | [optional] [default to null]
**Latitude** | **[]float32** | Location latitude, expressed in the range -90° to +90°. Cardinality greater than one only if \&quot;shape\&quot; equals 7. | [default to null]
**Longitude** | **[]float32** | Location longitude, expressed in the range -180° to +180°. Cardinality greater than one only if \&quot;shape\&quot; equals 7. | [default to null]
**OffsetAngle** | **int32** | Present only if \&quot;shape\&quot; equals 6 | [optional] [default to null]
**OrientationMajorAxis** | **int32** | Angle of orientation of the major axis, expressed in the range 0° to 180°, as defined in [14]. Present only if \&quot;shape\&quot; equals 4 or 6 | [optional] [default to null]
**Shape** | **int32** | Shape information, as detailed in [14], associated with the reported location coordinate: &lt;p&gt;1 &#x3D; ELLIPSOID_ARC &lt;p&gt;2 &#x3D; ELLIPSOID_POINT &lt;p&gt;3 &#x3D; ELLIPSOID_POINT_ALTITUDE &lt;p&gt;4 &#x3D; ELLIPSOID_POINT_ALTITUDE_UNCERT_ELLIPSOID &lt;p&gt;5 &#x3D; ELLIPSOID_POINT_UNCERT_CIRCLE &lt;p&gt;6 &#x3D; ELLIPSOID_POINT_UNCERT_ELLIPSE &lt;p&gt;7 &#x3D; POLYGON | [default to null]
**Timestamp** | [***TimeStamp**](TimeStamp.md) |  | [optional] [default to null]
**UncertaintyRadius** | **int32** | Present only if \&quot;shape\&quot; equals 6 | [optional] [default to null]
**Velocity** | [***LocationInfoVelocity**](LocationInfo_velocity.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


