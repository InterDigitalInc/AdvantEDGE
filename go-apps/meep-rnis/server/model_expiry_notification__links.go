/*
 * ETSI GS MEC 012 - Radio Network Information API
 *
 * The ETSI MEC ISG MEC012 Radio Network Information API described using OpenAPI.
 *
 * API version: 2.1.1
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package server

// List of hyperlinks related to the resource.
type ExpiryNotificationLinks struct {
	// Self referring URI. This shall be included in the response from the RNIS. The URI shall be unique within the RNI API as it acts as an ID for the subscription.
	Self string `json:"self"`
}
