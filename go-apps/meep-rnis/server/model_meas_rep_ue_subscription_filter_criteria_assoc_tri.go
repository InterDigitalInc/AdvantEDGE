/*
 * ETSI GS MEC 012 - Radio Network Information API
 *
 * The ETSI MEC ISG MEC012 Radio Network Information API described using OpenAPI.
 *
 * API version: 2.1.1
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package server

// List of filtering criteria for the subscription. Any filtering criteria from below, which is included in the request, shall also be included in the response.
type MeasRepUeSubscriptionFilterCriteriaAssocTri struct {
	// Unique identifier for the MEC application instance.
	AppInstanceId string `json:"appInstanceId,omitempty"`
	// 0 to N identifiers to associate the information for a specific UE or flow.
	AssociateId []AssociateId `json:"associateId,omitempty"`
	// E-UTRAN Cell Global Identifier.
	Ecgi []Ecgi `json:"ecgi,omitempty"`
	// Corresponds to a specific E-UTRAN UE Measurement Report trigger.
	Trigger []Trigger `json:"trigger,omitempty"`
}
